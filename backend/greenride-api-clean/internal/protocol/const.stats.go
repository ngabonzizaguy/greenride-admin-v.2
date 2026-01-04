package protocol

// 时区常量定义
const (
	// 中国时区
	TimezoneAsiaShanghai = "Asia/Shanghai"  // 上海（东八区，UTC+8）
	TimezoneAsiaHongKong = "Asia/Hong_Kong" // 香港（东八区，UTC+8）
	TimezoneAsiaTaipei   = "Asia/Taipei"    // 台北（东八区，UTC+8）
	TimezoneAsiaUrumqi   = "Asia/Urumqi"    // 乌鲁木齐（东六区，UTC+6）

	// 美国时区
	TimezoneAmericaNewYork    = "America/New_York"    // 纽约（东部时间，UTC-5/-4）
	TimezoneAmericaChicago    = "America/Chicago"     // 芝加哥（中部时间，UTC-6/-5）
	TimezoneAmericaDenver     = "America/Denver"      // 丹佛（山地时间，UTC-7/-6）
	TimezoneAmericaLosAngeles = "America/Los_Angeles" // 洛杉矶（太平洋时间，UTC-8/-7）
	TimezoneAmericaAnchorage  = "America/Anchorage"   // 安克雷奇（阿拉斯加时间，UTC-9/-8）
	TimezoneAmericaHonolulu   = "America/Honolulu"    // 火奴鲁鲁（夏威夷时间，UTC-10）

	// 欧洲时区
	TimezoneEuropeParis  = "Europe/Paris"  // 巴黎（中欧时间，UTC+1/+2）
	TimezoneEuropeLondon = "Europe/London" // 伦敦（格林威治时间，UTC+0/+1）
	TimezoneEuropeBerlin = "Europe/Berlin" // 柏林（中欧时间，UTC+1/+2）
	TimezoneEuropeMoscow = "Europe/Moscow" // 莫斯科（莫斯科时间，UTC+3）

	// 非洲时区
	TimezoneAfricaKigali       = "Africa/Kigali"       // 基加利，卢旺达（东非时间，UTC+2）
	TimezoneAfricaCairo        = "Africa/Cairo"        // 开罗，埃及（东欧时间，UTC+2）
	TimezoneAfricaLagos        = "Africa/Lagos"        // 拉各斯，尼日利亚（西非时间，UTC+1）
	TimezoneAfricaNairobi      = "Africa/Nairobi"      // 内罗毕，肯尼亚（东非时间，UTC+3）
	TimezoneAfricaJohannesburg = "Africa/Johannesburg" // 约翰内斯堡，南非（南非标准时间，UTC+2）

	// 卢旺达及周边国家时区
	TimezoneAfricaBujumbura   = "Africa/Bujumbura"     // 布琼布拉，布隆迪（东非时间，UTC+2）
	TimezoneAfricaKampala     = "Africa/Kampala"       // 坎帕拉，乌干达（东非时间，UTC+3）
	TimezoneAfricaDaresSalaam = "Africa/Dar_es_Salaam" // 达累斯萨拉姆，坦桑尼亚（东非时间，UTC+3）
	TimezoneAfricaKinshasa    = "Africa/Kinshasa"      // 金沙萨，刚果民主共和国（西非时间，UTC+1）
	TimezoneAfricaLubumbashi  = "Africa/Lubumbashi"    // 卢本巴希，刚果民主共和国（中非时间，UTC+2）

	// 其他常用时区
	TimezoneUTC             = "UTC"              // 世界协调时间（UTC+0）
	TimezoneAsiaTokyo       = "Asia/Tokyo"       // 东京，日本（日本标准时间，UTC+9）
	TimezoneAustraliaSydney = "Australia/Sydney" // 悉尼，澳大利亚（澳大利亚东部时间，UTC+10/+11）
	TimezoneAsiaKolkata     = "Asia/Kolkata"     // 加尔各答，印度（印度标准时间，UTC+5:30）
	TimezoneAsiaDubai       = "Asia/Dubai"       // 迪拜，阿联酋（海湾标准时间，UTC+4）
)

// 时间粒度类型常量
const (
	TimeTypeHour   = "hour"   // 小时粒度
	TimeTypeDay    = "day"    // 天粒度
	TimeTypeWeek   = "week"   // 周粒度
	TimeTypeMonth  = "month"  // 月粒度
	TimeTypeYear   = "year"   // 年粒度
	TimeTypeCustom = "custom" // 自定义时间段
)

// 时区组常量
var (
	// 中国时区组
	ChinaTimezones = []string{
		TimezoneAsiaShanghai,
		TimezoneAsiaHongKong,
		TimezoneAsiaTaipei,
		TimezoneAsiaUrumqi,
	}

	// 美国时区组
	USATimezones = []string{
		TimezoneAmericaNewYork,
		TimezoneAmericaChicago,
		TimezoneAmericaDenver,
		TimezoneAmericaLosAngeles,
		TimezoneAmericaAnchorage,
		TimezoneAmericaHonolulu,
	}

	// 欧洲时区组
	EuropeTimezones = []string{
		TimezoneEuropeParis,
		TimezoneEuropeLondon,
		TimezoneEuropeBerlin,
		TimezoneEuropeMoscow,
	}

	// 非洲时区组
	AfricaTimezones = []string{
		TimezoneAfricaKigali,
		TimezoneAfricaCairo,
		TimezoneAfricaLagos,
		TimezoneAfricaNairobi,
		TimezoneAfricaJohannesburg,
		TimezoneAfricaBujumbura,
		TimezoneAfricaKampala,
		TimezoneAfricaDaresSalaam,
		TimezoneAfricaKinshasa,
		TimezoneAfricaLubumbashi,
	}

	// 卢旺达及周边国家时区组
	RwandaRegionTimezones = []string{
		TimezoneAfricaKigali,
		TimezoneAfricaBujumbura,
		TimezoneAfricaKampala,
		TimezoneAfricaDaresSalaam,
		TimezoneAfricaLubumbashi,
	}

	// 全球主要时区组
	GlobalTimezones = []string{
		TimezoneUTC,
		TimezoneAsiaShanghai,
		TimezoneAmericaNewYork,
		TimezoneEuropeLondon,
		TimezoneAfricaKigali,
		TimezoneAsiaTokyo,
		TimezoneAustraliaSydney,
		TimezoneAsiaKolkata,
		TimezoneAsiaDubai,
	}
)

// 默认时区（应用默认使用卢旺达时区）
const DefaultTimezone = TimezoneAfricaKigali
