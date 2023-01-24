package lib

import (
	"fmt"
	"net/http"
	"path/filepath"
	"reflect"
	"strings"
	"sync"

	"go/goserver.io/constants"
	"go/goserver.io/utils"

	"github.com/dgofman/pongo2"
)

/**
 * Local type definitions
 */
type commondNode struct {
	content string
	tagName string
}

type Form struct {
	r    *http.Request
	dict map[string]interface{}
	list []*Field
}

type Field struct {
	id        string
	attrs     string
	Value     *string      `json:"value"`
	Auto_id   string       `json:"auto_id"`
	Label_tag *labelExtend `json:"label_tag"`
	Field     *fieldExtend `json:"field"`
	Errors    []string     `json:"errors"`
}

type labelExtend struct {
	id    string
	label string
}

type fieldExtend struct {
	Required bool              `json:"required"`
	Widget   map[string]string `json:"widget"`
}

/**
 * Represents Label attributes as HTML string
 */
func (l *labelExtend) String() string {
	return fmt.Sprintf("<label for=\"id_%s\">%s</label>", l.id, l.label)
}

/**
 * Represents Field attributes as HTML string
 */
func (f *Field) String() string {
	attrs := f.attrs
	if f.Value != nil {
		attrs += fmt.Sprintf(` value="%s"`, *f.Value)
	}
	return fmt.Sprintf(`<input type="%s" name="%s" %s id="id_%s">`, f.Field.Widget["input_type"], f.id, attrs, f.id)
}

/**
 * Add error objects to the Form
 */
func (f *Field) AddError(err string) {
	f.Errors = append(f.Errors, err)
}

/**
 * Dynamic add/update the field value
 */
func (f *Field) SetValue(val string) {
	f.Value = &val
}

/**
 * Apply Django "static" tag parser
 */
func staticParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &commondNode{tagName: "static"}
	pathToken := arguments.MatchType(pongo2.TokenString)
	if pathToken != nil {
		node.content = constants.URL_Static + pathToken.Val
	}
	return node, nil
}

/**
 * Apply Django "url" tag parser
 */
func urlParser(doc *pongo2.Parser, start *pongo2.Token, arguments *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
	node := &commondNode{tagName: "url"}
	urlToken := arguments.MatchType(pongo2.TokenString)
	if urlToken != nil {
		node.content = urlToken.Val
	}
	return node, nil
}

/**
 * Replaced "url" tag in the HTML by form parameter starts with "_url_"
 */
func (node *commondNode) Execute(ctx *pongo2.ExecutionContext, writer pongo2.TemplateWriter) *pongo2.Error {
	if node.tagName == "url" {
		forCtx := pongo2.NewChildExecutionContext(ctx)
		value := forCtx.Public["_url_"+node.content]
		if value == nil {
			utils.LogError("Missing context parameter: _url_" + node.content)
		} else {
			writer.WriteString(value.(string))
		}
	} else if node.content != "" {
		writer.WriteString(node.content)
	}
	return nil
}

/**
 * Return absolute path to the .../templates directory
 */
func TemplatePath() string {
	return utils.Abs(filepath.Join(constants.Props.TEMPLATE_DIR))
}

var once sync.Once

/**
 * Parse and create a Template object using file name
 */
func Template(filename string) *pongo2.Template {
	once.Do(func() {
		pongo2.RegisterTag("load", func(*pongo2.Parser, *pongo2.Token, *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
			return &commondNode{}, nil
		})
		pongo2.RegisterTag("csrf_token", func(*pongo2.Parser, *pongo2.Token, *pongo2.Parser) (pongo2.INodeTag, *pongo2.Error) {
			return &commondNode{}, nil
		})
		pongo2.RegisterTag("static", staticParser)
		pongo2.RegisterTag("url", urlParser)
	})
	tplSet := pongo2.NewSet("template-set", pongo2.MustNewLocalFileSystemLoader(TemplatePath()))
	t, err := tplSet.FromFile(filename)
	if err != nil {
		panic(err)
	}
	return t
}

/**
 * Build Template
 */
func BuildTemplate(w http.ResponseWriter, filename string, context pongo2.Context) {
	ExecTemplate(w, Template(filename), context)
}

/**
 * Execute Template and pass runtime arguments
 */
func ExecTemplate(w http.ResponseWriter, t *pongo2.Template, context pongo2.Context) {
	header := w.Header()
	header.Set("X-Frame-Options", "SAMEORIGIN")
	header.Set("Cache-Control", "no-cache, private, max-age=0")
	err := t.ExecuteWriter(context, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type GetFieldsFunc func(ctx *pongo2.ExecutionContext) interface{}

/**
 * Return array or map of Field objects based on Django "form.field" expression
 */
func (f *Form) GetFields(ctx *pongo2.ExecutionContext) interface{} {
	isLoop := ctx.Private["forloop"]
	if isLoop != nil {
		el := reflect.ValueOf(isLoop).Elem()
		tag := el.FieldByName("TagNode").Elem()
		key := tag.FieldByName("key").String()
		if key == "field" {
			return f.list
		}
	}
	return f.dict
}

/**
 * Defined the default template-form arguments
 */
func CreateContent(r *http.Request) pongo2.Context {
	urlencode := map[string]string{
		"urlencode": r.URL.RawQuery,
	}

	if constants.Props.VIEW == "light" {
		urlencode["urlencode"] += "&view=light"
	}

	meta := map[string]interface{}{}
	if r.Header.Get("X-Auth-User-Fn") != "" {
		meta["HTTP_X_AUTH_USER_FN"] = r.Header.Get("X-Auth-User-Fn")
	}
	if r.Header.Get("X-Auth-User-Avatar") != "" {
		meta["HTTP_X_AUTH_USER_AVATAR"] = r.Header.Get("X-Auth-User-Avatar")
	}

	ctx := make(pongo2.Context)
	ctx["_url_device_info"] = constants.URL_DeviceInfo
	ctx["_url_ajax_ipdata"] = constants.URL_AjaxIpdata
	ctx["_url_ajax_metrics"] = constants.URL_AjaxMetrics
	ctx["_url_ajax_verify_email"] = constants.URL_VerifyEmail
	ctx["_url_add_client"] = constants.URL_AddClient
	ctx["request"] = map[string]interface{}{
		"GET":  urlencode,
		"POST": urlencode,
		"META": meta,
	}
	return ctx
}

/**
 * Initialize a Form object
 */
func InitForm(r *http.Request, dict map[string]interface{}) *Form {
	if dict == nil {
		dict = make(map[string]interface{})
	}
	return &Form{
		r:    r,
		dict: dict,
		list: make([]*Field, 0),
	}
}

/**
 * Add or Update the form values
 */
func (form *Form) Set(key string, value interface{}) {
	form.dict[key] = value
}

/**
 * Get the form value or field
 */
func (form *Form) Get(key string) interface{} {
	return form.dict[key]
}

/**
 * Convert POST form request to JSON
 */
func (form *Form) ToMap() map[string]interface{} {
	return ToMap(form.r)
}

/**
 * Convert POST form request to JSON
 */
func ToMap(r *http.Request) map[string]interface{} {
	dict := map[string]interface{}{}
	for key, value := range r.PostForm {
		dict[key] = value[0]
	}
	return dict
}

/**
 * Create a Django Field object
 */
func (form *Form) CharField(id string, label string, input_type string, attrs string, value ...interface{}) {
	field := &Field{
		id:      id,
		attrs:   attrs,
		Value:   nil,
		Auto_id: fmt.Sprintf(`id_%s`, id),
		Label_tag: &labelExtend{
			id:    id,
			label: label,
		},
		Field: &fieldExtend{
			Required: strings.Contains(attrs, "required"),
			Widget: map[string]string{
				"input_type": input_type,
			},
		},
		Errors: make([]string, 0),
	}
	if len(value) == 0 {
		field.SetValue(form.r.FormValue(id))
	} else if value[0] != nil {
		field.SetValue(value[0].(string))
	}
	form.dict[id] = field
	form.list = append(form.list, field)
}
