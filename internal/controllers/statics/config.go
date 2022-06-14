package statics

type (
	Config struct {
		Content      ContentItem `yaml:"content"`
		CDN          CDNItem     `yaml:"cdn"`
		AllowDomains []string    `yaml:"allow_domains"`
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
	v.AllowDomains = []string{"127.0.0.1"}
	v.CDN = CDNItem{
		Domain: "http://127.0.0.1",
		Path:   "/tmp/solocms/cdn",
	}
	v.Content = ContentItem{
		Path: "/tmp/solocms/data",
	}
}
