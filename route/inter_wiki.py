from .tool.func import *

def inter_wiki(conn, tool):
    curs = conn.cursor()

    div = '<table id="main_table_set">'
    div += '<tr id="main_table_top_tr">'
    
    div += '<td id="main_table_width">A</td>'
    div += '<td id="main_table_width">B</td>'
    div += '<td id="main_table_width">C</td>'
    
    div += '</tr>'
    
    admin = admin_check()

    if tool == 'inter_wiki':
        plus_link = 'plus_inter_wiki'
        title = load_lang('interwiki_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'inter_wiki'"))
    elif tool == 'email_filter':
        plus_link = 'plus_email_filter'
        title = load_lang('email_filter_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'email'"))
    elif tool == 'name_filter':
        plus_link = 'plus_name_filter'
        title = load_lang('id_filter_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'name'"))
    elif tool == 'edit_filter':
        plus_link = 'plus_edit_filter'
        title = load_lang('edit_filter_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'regex_filter'"))
    elif tool == 'file_filter':
        plus_link = 'plus_file_filter'
        title = load_lang('file_filter_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'file'"))
    elif tool == 'file_filter':
        plus_link = 'plus_file_filter'
        title = load_lang('file_filter_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'file'"))
    elif tool == 'image_license':
        plus_link = 'plus_image_license'
        title = load_lang('image_license_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'image_license'"))
    elif tool == 'extension_filter':
        plus_link = 'plus_extension_filter'
        title = load_lang('extension_filter_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'extension'"))
    else:
        plus_link = 'plus_edit_top'
        title = load_lang('edit_tool_list')

        curs.execute(db_change("select html, plus, plus_t from html_filter where kind = 'edit_top'"))

    db_data = curs.fetchall()
    for data in db_data:
        div += '<tr>'
        div += '<td>'

        div += data[0]
        if admin == 1:
            div += ' <a href="/' + tool + '/add/' + url_pas(data[0]) + '">(' + load_lang('edit') + ')</a>'
            div += ' <a href="/' + tool + '/del/' + url_pas(data[0]) + '">(' + load_lang('delete') + ')</a>'
        
        div += '</td>'

        if tool == 'inter_wiki':
            div += '<td><a id="out_link" href="' + data[1] + '">' + data[1] + '</a></td>'
        else:
            div += '<td>' + data[1] + '</td>'
            
        div += '<td>' + data[2] + '</td>'
        div += '</tr>'
        
    div += '</table>'
            
    if admin == 1:
        div += '<hr class="main_hr">'
        div += '<a href="/' + tool + '/add">(' + load_lang('add') + ')</a>'

    return easy_minify(flask.render_template(skin_check(),
        imp = [title, wiki_set(), wiki_custom(), wiki_css([0, 0])],
        data = div,
        menu = [['manager/1', load_lang('return')]]
    ))