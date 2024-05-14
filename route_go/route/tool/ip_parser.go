package tool

import (
	"database/sql"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func IP_or_user(ip string) bool {
	match, _ := regexp.MatchString("(\\.|:)", ip)
	if match {
		return true
	} else {
		return false
	}
}

func Get_level(db *sql.DB, db_set map[string]string, ip string) []string {
	var level string
	var exp string
	var max_exp string

	stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'level'"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(ip).Scan(&level)
	if err != nil {
		if err == sql.ErrNoRows {
			level = "0"
		} else {
			log.Fatal(err)
		}
	}

	stmt, err = db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'experience'"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	err = stmt.QueryRow(ip).Scan(&exp)
	if err != nil {
		if err == sql.ErrNoRows {
			exp = "0"
		} else {
			log.Fatal(err)
		}
	}

	level_int, _ := strconv.Atoi(level)
	max_exp = strconv.Itoa(level_int*50 + 500)

	return []string{level, exp, max_exp}
}

func Get_user_auth(db *sql.DB, db_set map[string]string, ip string) string {
	if !IP_or_user(ip) {
		var auth string

		stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'acl'"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		err = stmt.QueryRow(ip).Scan(&auth)
		if err != nil {
			if err == sql.ErrNoRows {
				auth = "user"
			} else {
				log.Fatal(err)
			}
		}

		if auth != "user" && auth != "ban" {
			return auth
		} else {
			return ""
		}
	}

	return ""
}

func Get_auth_group_info(db *sql.DB, db_set map[string]string, auth string) map[string]bool {
	stmt, err := db.Prepare(DB_change(db_set, "select name from alist where name = ?"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(auth)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	data_list := map[string]bool{}

	for rows.Next() {
		var name string

		err := rows.Scan(&name)
		if err != nil {
			log.Fatal(err)
		}

		data_list[name] = true
	}

	return data_list
}

func IP_preprocess(db *sql.DB, db_set map[string]string, ip string, my_ip string) []string {
	var ip_view string
	var user_name_view string

	ip_split := strings.Split(ip, ":")
	if len(ip_split) != 1 && ip_split[0] == "tool" {
		return []string{ip, ""}
	}

	err := db.QueryRow(DB_change(db_set, "select data from other where name = 'ip_view'")).Scan(&ip_view)
	if err != nil {
		if err == sql.ErrNoRows {
			ip_view = ""
		} else {
			log.Fatal(err)
		}
	}

	err = db.QueryRow(DB_change(db_set, "select data from other where name = 'user_name_view'")).Scan(&user_name_view)
	if err != nil {
		if err == sql.ErrNoRows {
			user_name_view = ""
		} else {
			log.Fatal(err)
		}
	}

	if Get_user_auth(db, db_set, my_ip) != "" {
		ip_view = ""
		user_name_view = ""
	}

	ip_change := ""
	if IP_or_user(ip) {
		if ip_view != "" && ip != my_ip {
			hash_ip := Sha224(ip)
			ip = hash_ip[:10]
			ip_change = "true"
		}
	} else {
		if user_name_view != "" {
			var sub_user_name string

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where id = ? and name = 'sub_user_name'"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(ip).Scan(&sub_user_name)
			if err != nil {
				if err == sql.ErrNoRows {
					sub_user_name = Get_language(db, db_set, "member", false)
				} else {
					log.Fatal(err)
				}
			}

			if sub_user_name == "" {
				sub_user_name = Get_language(db, db_set, "member", false)
			}

			ip = sub_user_name
			ip_change = "true"
		} else {
			var user_name string

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where name = 'user_name' and id = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(ip).Scan(&user_name)
			if err != nil {
				if err == sql.ErrNoRows {
					user_name = ip
				} else {
					log.Fatal(err)
				}
			}

			if user_name == "" {
				user_name = ip
			}

			ip = user_name
		}
	}

	return []string{ip, ip_change}
}

func IP_menu(db *sql.DB, db_set map[string]string, ip string, my_ip string, option string) map[string][][]string {
	menu := map[string][][]string{}

	if ip == my_ip && option == "" {
		stmt, err := db.Prepare(DB_change(db_set, "select count(*) from user_notice where name = ? and readme = ''"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var alarm_count string

		err = stmt.QueryRow(my_ip).Scan(&alarm_count)
		if err != nil {
			if err == sql.ErrNoRows {
				alarm_count = "0"
			} else {
				log.Fatal(err)
			}
		}

		if IP_or_user(my_ip) {
			menu[Get_language(db, db_set, "login", false)] = [][]string{
				{"/login", Get_language(db, db_set, "login", false)},
				{"/register", Get_language(db, db_set, "register", false)},
				{"/change", Get_language(db, db_set, "user_setting", false)},
				{"/login/find", Get_language(db, db_set, "password_search", false)},
				{"/alarm" + Url_parser(my_ip), Get_language(db, db_set, "alarm", false) + " (" + alarm_count + ")"},
			}
		} else {
			menu[Get_language(db, db_set, "login", false)] = [][]string{
				{"/logout", Get_language(db, db_set, "logout", false)},
				{"/change", Get_language(db, db_set, "user_setting", false)},
			}

			menu[Get_language(db, db_set, "tool", false)] = [][]string{
				{"/watch_list", Get_language(db, db_set, "watch_list", false)},
				{"/star_doc", Get_language(db, db_set, "star_doc", false)},
				{"/challenge", Get_language(db, db_set, "challenge_and_level_manage", false)},
				{"/acl/user:" + Url_parser(my_ip), Get_language(db, db_set, "user_document_acl", false)},
				{"/alarm" + Url_parser(my_ip), Get_language(db, db_set, "alarm", false) + " (" + alarm_count + ")"},
			}
		}
	}

	auth_name := Get_user_auth(db, db_set, my_ip)
	if auth_name != "" {
		menu[Get_language(db, db_set, "admin", false)] = [][]string{
			{"/auth/give/ban/" + Url_parser(ip), Get_language(db, db_set, "ban", false) + " | " + Get_language(db, db_set, "release", false)},
			{"/list/user/check/" + Url_parser(ip), Get_language(db, db_set, "check", false)},
		}
	}

	menu[Get_language(db, db_set, "other", false)] = [][]string{
		{"/record/" + Url_parser(ip), Get_language(db, db_set, "edit_record", false)},
		{"/record/topic/" + Url_parser(ip), Get_language(db, db_set, "discussion_record", false)},
		{"/record/bbs/" + Url_parser(ip), Get_language(db, db_set, "bbs_record", false)},
		{"/record/bbs_comment/" + Url_parser(ip), Get_language(db, db_set, "bbs_comment_record", false)},
		{"/topic/user:" + Url_parser(ip), Get_language(db, db_set, "user_discussion", false)},
		{"/count/" + Url_parser(ip), Get_language(db, db_set, "count", false)},
	}

	return menu
}

func IP_parser(db *sql.DB, db_set map[string]string, ip string, my_ip string) string {
	ip_pre_data := IP_preprocess(db, db_set, ip, my_ip)
	if ip_pre_data[0] == "" {
		return ""
	}

	if ip_pre_data[1] != "" {
		return ip_pre_data[0]
	} else {
		raw_ip := ip
		ip = HTML_escape(ip_pre_data[0])

		if !IP_or_user(raw_ip) {
			var user_name_level string
			var user_title string

			err := db.QueryRow(DB_change(db_set, "select data from other where name = 'user_name_level'")).Scan(&user_name_level)
			if err != nil {
				if err == sql.ErrNoRows {
					user_name_level = ""
				} else {
					log.Fatal(err)
				}
			}

			if user_name_level != "" {
				level_data := Get_level(db, db_set, raw_ip)
				ip += "<sup>" + level_data[0] + "</sup>"
			}

			ip = "<a href=\"/w/" + Url_parser("user:"+raw_ip) + "\">" + ip + "</a>"

			stmt, err := db.Prepare(DB_change(db_set, "select data from user_set where name = 'user_title' and id = ?"))
			if err != nil {
				log.Fatal(err)
			}
			defer stmt.Close()

			err = stmt.QueryRow(raw_ip).Scan(&user_title)
			if err != nil {
				if err == sql.ErrNoRows {
					user_title = ""
				} else {
					log.Fatal(err)
				}
			}

			if Get_user_auth(db, db_set, raw_ip) != "" {
				ip = "<b>" + ip + "</b>"
			}

			ip = user_title + ip
		}

		ip += " <a href=\"javascript:opennamu_do_ip_click(this);\">(" + Get_language(db, db_set, "tool", false) + ")</a>"

		return ip
	}
}
