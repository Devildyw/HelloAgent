package model

// LocationResponse 和风天气城市查询API响应
type LocationResponse struct {
	Code     string     `json:"code"`
	Location []Location `json:"location"`
	Refer    Refer      `json:"refer"`
}

// Location 城市位置信息
type Location struct {
	Name      string `json:"name"`      // 城市名称
	ID        string `json:"id"`        // 城市ID
	Lat       string `json:"lat"`       // 纬度
	Lon       string `json:"lon"`       // 经度
	Adm2      string `json:"adm2"`      // 上级行政区划名称
	Adm1      string `json:"adm1"`      // 一级行政区划名称
	Country   string `json:"country"`   // 国家名称
	Tz        string `json:"tz"`        // 时区
	UtcOffset string `json:"utcOffset"` // UTC偏移
	IsDst     string `json:"isDst"`     // 是否为夏令时 0=否 1=是
	Type      string `json:"type"`      // 地区类型
	Rank      string `json:"rank"`      // 地区评分
	FxLink    string `json:"fxLink"`    // 该地区的天气预报网页链接
}

// Refer 数据来源信息
type Refer struct {
	Sources []string `json:"sources"` // 数据来源
	License []string `json:"license"` // 许可证
}

// WeatherResponse 和风天气实时天气查询API响应
type WeatherResponse struct {
	Code       string     `json:"code"`       // 状态码
	UpdateTime string     `json:"updateTime"` // 当前API的最近更新时间
	FxLink     string     `json:"fxLink"`     // 当前数据的响应式页面
	Now        WeatherNow `json:"now"`        // 实时天气数据
	Refer      Refer      `json:"refer"`      // 数据来源信息
}

// WeatherNow 实时天气数据
type WeatherNow struct {
	ObsTime   string `json:"obsTime"`   // 数据观测时间
	Temp      string `json:"temp"`      // 温度，默认单位：摄氏度
	FeelsLike string `json:"feelsLike"` // 体感温度，默认单位：摄氏度
	Icon      string `json:"icon"`      // 天气状况图标代码
	Text      string `json:"text"`      // 天气状况的文字描述
	Wind360   string `json:"wind360"`   // 风向360角度
	WindDir   string `json:"windDir"`   // 风向
	WindScale string `json:"windScale"` // 风力等级
	WindSpeed string `json:"windSpeed"` // 风速，公里/小时
	Humidity  string `json:"humidity"`  // 相对湿度，百分比数值
	Precip    string `json:"precip"`    // 当前小时累计降水量，默认单位：毫米
	Pressure  string `json:"pressure"`  // 大气压强，默认单位：百帕
	Vis       string `json:"vis"`       // 能见度，默认单位：公里
	Cloud     string `json:"cloud"`     // 云量，百分比数值
	Dew       string `json:"dew"`       // 露点温度
}

// SearchResponse 搜索API响应
type SearchResponse struct {
	Query             string         `json:"query"`               // 搜索关键词
	FollowUpQuestions []string       `json:"follow_up_questions"` // 后续推荐问题
	Answer            *string        `json:"answer"`              // 答案摘要（可能为null）
	Images            []string       `json:"images"`              // 相关图片URL列表
	Results           []SearchResult `json:"results"`             // 搜索结果列表
	ResponseTime      float64        `json:"response_time"`       // 响应时间（秒）
	RequestID         string         `json:"request_id"`          // 请求唯一标识
}

// SearchResult 单条搜索结果
type SearchResult struct {
	URL        string  `json:"url"`         // 结果页面URL
	Title      string  `json:"title"`       // 页面标题
	Content    string  `json:"content"`     // 内容摘要
	Score      float64 `json:"score"`       // 相关性评分
	RawContent *string `json:"raw_content"` // 原始内容（可能为null）
}
