package controllers

type (
	Config struct {
		Content      ContentItem       `yaml:"content"`
		CDN          CDNItem           `yaml:"cdn"`
		AllowDomains []string          `yaml:"allow_domains"`
		Langs        []string          `yaml:"langs"`
		AdminAuth    map[string]string `yaml:"admin_auth"`
	}

	ContentItem struct {
		Path string `yaml:"path"`
	}

	CDNItem struct {
		Domain string `yaml:"domain"`
		Path   string `yaml:"path"`
	}
)

func (v *Config) Default() {
	v.Langs = []string{"en"}
	v.AllowDomains = []string{"127.0.0.1"}
	v.CDN = CDNItem{
		Domain: "http://127.0.0.1",
		Path:   "/tmp/solocms/cdn",
	}
	v.Content = ContentItem{
		Path: "/tmp/solocms/data",
	}
}

func (v *Config) HasLang(lang string) bool {
	for _, s := range v.Langs {
		if s == lang {
			return true
		}
	}
	return false
}
