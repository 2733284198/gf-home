package a_document

import (
    "github.com/gogf/gf-home/app/service/document"
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/os/gfile"
    "net/http"
)

// 文档首页
func Index(r *ghttp.Request) {
    if r.IsAjaxRequest() {
        serveMarkdownAjax(r)
        return
    }
    path := r.Get("path")
    if path == "" {
        if r.URL.RawQuery != "" {
            r.Response.RedirectTo("/index?" + r.URL.RawQuery)
        } else {
            r.Response.RedirectTo("/index")
        }
        return
    }
    config := g.Config()
    // 如果是静态文件请求，那么表示Web Server没有找到该文件，那么直接404，本接口不支持待后缀的静态文件处理。
    // 由于路由规则比较宽，这里也会有未存在的静态文件请求匹配进来。
    if gfile.Ext(path) != "" {
        r.Response.WriteStatus(http.StatusNotFound)
        return
    }
    // 菜单内容
    baseTitle := config.GetString("document.title")
    title     := s_document.GetTitleByPath(path)
    if title == "" {
        title = "404 NOT FOUND"
    }
    title += " - " + config.GetString("document.title")
    // markdown内容
    mdMainContent       := s_document.GetMarkdown(path)
    mdMainContentParsed := s_document.ParseMarkdown(mdMainContent)
    r.Response.WriteTpl("document/index.html", g.Map {
        "title"               : title,
        "baseTitle"           : baseTitle,
        "mdMenuContentParsed" : s_document.GetParsed("menus"),
        "mdMainContentParsed" : mdMainContentParsed,
        "mdMainContent"       : mdMainContent,
    })
}

// 文档更新hook
func UpdateHook(r *ghttp.Request) {
    if r.Get("password") == g.Config().GetString("document.hook") {
        s_document.UpdateDocGit()
        r.Response.Write("ok")
    } else {
        r.Response.WriteStatus(443)
    }
}

// 搜索文档
func Search(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code" : 1,
        "msg"  : "",
        "data" : s_document.SearchMdByKey(r.GetString("key")),
    })
}

// 处理ajax请求
func serveMarkdownAjax(r *ghttp.Request) {
    r.Response.WriteJson(g.Map{
        "code" : 1,
        "msg"  : "",
        "data" : s_document.GetMarkdown(r.Get("path", "index")),
    })
}