package schema

import (
	"time"

	"gin-casbin/pkg/util"
)

// Tenant 租户对象
type Tenant struct {
	ID                string    `json:"id"`                      // 唯一标识
	Name              string    `json:"name" binding:"required"` // 租户名称
	URL               string    `json:"url" validate:"url"`      // 租户URL
	LogoURL           string    `json:"logo_url" validate:"url"` // 租户LOGO URL
	Timezone          string    `json:"timezone"`                // 时区
	Language          string    `json:"language"`                // 语言
	Theme             string    `json:"theme"`                   // 默认主题
	Description       string    `json:"description"`             // 描述
	Phone             string    `json:"phone"`                   // 电话
	Details           string    `json:"details"`                 // 详细
	Status            int       `json:"status"`                  // 用户状态(1:启用 2:停用)
	Creator           string    `json:"creator"`                 // 创建者
	CreatedAt         time.Time `json:"created_at"`              // 创建时间
	UpdatedAt         time.Time `json:"updated_at"`              // 更新时间
	Address           Address   `json:"address"`                 // 租户地址
	Administrator     *User     `json:"administrator"`           // 租户管理员
	MaxQrQty          int64     `json:"max_qr_qty"`              // 已购买QR数量
	UsedQrQty         int64     `json:"used_qr_qty"`             // 已使用QR数量
	TotalOrderQty     int64     `json:"total_order_qty"`         // 已提交订单数
	ProcessedOrderQty int64     `json:"processed_order_qty"`     // 已处理QR数量
}

func (a *Tenant) String() string {
	return util.JSONMarshalToString(a)
}

// TenantQueryParam 查询条件
type TenantQueryParam struct {
	QrCodeID string   `form:"qr_code_id"`
	ID       string   `form:"id"`
	IDs      []string `form:"ids"`
	PaginationParam
}

// TenantQueryOptions 查询可选参数项
type TenantQueryOptions struct {
	OrderFields []*OrderField // 排序字段
}

// TenantGetOptions Get查询可选参数项
type TenantGetOptions struct {
}

// TenantQueryResult 查询结果
type TenantQueryResult struct {
	Data       Tenants
	PageResult *PaginationResult
}

// Tenants 租户列表
type Tenants []*Tenant

// ToMap 转换为键值存储
func (a Tenants) ToMap() map[string]*Tenant {
	m := make(map[string]*Tenant)
	for _, item := range a {
		m[item.ID] = item
	}
	return m
}

// ToUserIDMap 转换为用户ID映射
func (a UserTenants) ToUserIDMap() map[string]UserTenant {
	m := make(map[string]UserTenant)
	for _, item := range a {
		m[item.UserID] = *item
	}
	return m
}

// // Tenant 租户输入
// type TenantInput struct {
// 	Name          string    `form:"name" binding:"required"`     // 租户名称
// 	URL           string    `form:"url" validate:"url"`          // 租户URL
// 	Timezone      string    `form:"timezone" binding:"required"` // 时区
// 	Language      string    `form:"language" binding:"required"` // 语言
// 	Theme         string    `form:"theme" binding:"required"`    // 默认主题
// 	Description   string    `form:"description"`                 // 描述
// 	Details       string    `form:"details"`                     // 详细
// 	Creator       string    `form:"creator"`                     // 创建者
// 	CreatedAt     time.Time `form:"created_at"`                  // 创建时间
// 	UpdatedAt     time.Time `form:"updated_at"`                  // 更新时间
// 	Address       Address   `form:"address"`                     // 租户地址
// 	Administrator User      `form:"administrator"`               // 租户管理员
// }

// // validator/baked_in.go
// bakedInValidators = map[string]Func{
// 	"required":         hasValue,
// 	"isdefault":        isDefault,
// 	"len":              hasLengthOf,
// 	"min":              hasMinOf,
// 	"max":              hasMaxOf,
// 	"eq":               isEq,
// 	"ne":               isNe,
// 	"lt":               isLt,
// 	"lte":              isLte,
// 	"gt":               isGt,
// 	"gte":              isGte,
// 	"eqfield":          isEqField,
// 	"eqcsfield":        isEqCrossStructField,
// 	"necsfield":        isNeCrossStructField,
// 	"gtcsfield":        isGtCrossStructField,
// 	"gtecsfield":       isGteCrossStructField,
// 	"ltcsfield":        isLtCrossStructField,
// 	"ltecsfield":       isLteCrossStructField,
// 	"nefield":          isNeField,
// 	"gtefield":         isGteField,
// 	"gtfield":          isGtField,
// 	"ltefield":         isLteField,
// 	"ltfield":          isLtField,
// 	"alpha":            isAlpha,
// 	"alphanum":         isAlphanum,
// 	"alphaunicode":     isAlphaUnicode,
// 	"alphanumunicode":  isAlphanumUnicode,
// 	"numeric":          isNumeric,
// 	"number":           isNumber,
// 	"hexadecimal":      isHexadecimal,
// 	"hexcolor":         isHEXColor,
// 	"rgb":              isRGB,
// 	"rgba":             isRGBA,
// 	"hsl":              isHSL,
// 	"hsla":             isHSLA,
// 	"email":            isEmail,
// 	"url":              isURL,
// 	"uri":              isURI,
// 	"file":             isFile,
// 	"base64":           isBase64,
// 	"base64url":        isBase64URL,
// 	"contains":         contains,
// 	"containsany":      containsAny,
// 	"containsrune":     containsRune,
// 	"excludes":         excludes,
// 	"excludesall":      excludesAll,
// 	"excludesrune":     excludesRune,
// 	"isbn":             isISBN,
// 	"isbn10":           isISBN10,
// 	"isbn13":           isISBN13,
// 	"eth_addr":         isEthereumAddress,
// 	"btc_addr":         isBitcoinAddress,
// 	"btc_addr_bech32":  isBitcoinBech32Address,
// 	"uuid":             isUUID,
// 	"uuid3":            isUUID3,
// 	"uuid4":            isUUID4,
// 	"uuid5":            isUUID5,
// 	"ascii":            isASCII,
// 	"printascii":       isPrintableASCII,
// 	"multibyte":        hasMultiByteCharacter,
// 	"datauri":          isDataURI,
// 	"latitude":         isLatitude,
// 	"longitude":        isLongitude,
// 	"ssn":              isSSN,
// 	"ipv4":             isIPv4,
// 	"ipv6":             isIPv6,
// 	"ip":               isIP,
// 	"cidrv4":           isCIDRv4,
// 	"cidrv6":           isCIDRv6,
// 	"cidr":             isCIDR,
// 	"tcp4_addr":        isTCP4AddrResolvable,
// 	"tcp6_addr":        isTCP6AddrResolvable,
// 	"tcp_addr":         isTCPAddrResolvable,
// 	"udp4_addr":        isUDP4AddrResolvable,
// 	"udp6_addr":        isUDP6AddrResolvable,
// 	"udp_addr":         isUDPAddrResolvable,
// 	"ip4_addr":         isIP4AddrResolvable,
// 	"ip6_addr":         isIP6AddrResolvable,
// 	"ip_addr":          isIPAddrResolvable,
// 	"unix_addr":        isUnixAddrResolvable,
// 	"mac":              isMAC,
// 	"hostname":         isHostnameRFC952,  // RFC 952
// 	"hostname_rfc1123": isHostnameRFC1123, // RFC 1123
// 	"fqdn":             isFQDN,
// 	"unique":           isUnique,
// 	"oneof":            isOneOf,
// 	"html":             isHTML,
// 	"html_encoded":     isHTMLEncoded,
// 	"url_encoded":      isURLEncoded,
// }
