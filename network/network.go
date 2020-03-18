package network

import (
	"github.com/go-ini/ini"
)

const (
	// 配置路径
	configPath = "/mnt/yaffs2/net.conf"

	// 默认配置
	defaultIP        = "192.168.1.233"
	defaultMask      = "255.255.255.0"
	defaultGatewayIP = "192.168.1.1"
)

// Config 配置文件网络配置参数
type Config struct {
	Auto         bool   `json:"auto"`                                             // 为true忽略其它参数,采用动态ip
	IPAddress    string `json:"ipAddress" validate:"ipv4"`                        // 静态Ip地址
	NetMask      string `json:"netMask" validate:"ipv4"`                          // 掩码
	Gateway      string `json:"gateway" validate:"ipv4"`                          // 网关地址
	PrimaryDNS   string `json:"primaryDNS,omitempty" validate:"ipv4|isdefault"`   // 主DNS服务器
	SecondaryDNS string `json:"secondaryDNS,omitempty" validate:"ipv4|isdefault"` // 次DNS服务器
	Mac          string `json:"-"`                                                // 这个是不准改的
}

func init() {
	ini.PrettyFormat = false
}

// GetConfig 获得配置文件的相关配置
func GetConfig() Config {
	cfgFile, err := ini.Load(configPath)
	if err != nil {
		return Config{}
	}
	sec := cfgFile.Section(ini.DefaultSection)
	return Config{
		sec.Key("AUTO").MustBool(false),
		sec.Key("IPADDR").Value(),
		sec.Key("NETMASK").Value(),
		sec.Key("GATEWAY").Value(),
		sec.Key("PrimaryDNS").Value(),
		sec.Key("SecondaryDNS").Value(),
		sec.Key("MAC").Value(),
	}
}

// SaveConfig 保存配置
func SaveConfig(newCfg *Config) error {
	cfgFile, err := ini.Load(configPath)
	if err != nil {
		return err
	}
	// NOTE: mac 不变
	sec := cfgFile.Section(ini.DefaultSection)
	if newCfg.Auto {
		sec.Key("AUTO").SetValue("true")
	} else {
		sec.DeleteKey("AUTO")
		sec.Key("IPADDR").SetValue(newCfg.IPAddress)
		sec.Key("NETMASK").SetValue(newCfg.NetMask)
		sec.Key("GATEWAY").SetValue(newCfg.Gateway)
		if newCfg.PrimaryDNS == "" {
			sec.DeleteKey("PrimaryDNS")
		} else {
			sec.Key("PrimaryDNS").SetValue(newCfg.PrimaryDNS)
		}
		if newCfg.SecondaryDNS == "" {
			sec.DeleteKey("SecondaryDNS")
		} else {
			sec.Key("SecondaryDNS").SetValue(newCfg.SecondaryDNS)
		}
	}
	return cfgFile.SaveTo(configPath)
}

// FactoryConfig 恢复出厂设置
func FactoryConfig() error {
	cfgFile, err := ini.Load(configPath)
	if err != nil {
		return err
	}

	// NOTE: mac 不变
	sec := cfgFile.Section(ini.DefaultSection)
	sec.DeleteKey("AUTO")
	sec.Key("IPADDR").SetValue(defaultIP)
	sec.Key("NETMASK").SetValue(defaultMask)
	sec.Key("GATEWAY").SetValue(defaultGatewayIP)
	sec.DeleteKey("PrimaryDNS")
	sec.DeleteKey("SecondaryDNS")
	return cfgFile.SaveTo(configPath)
}
