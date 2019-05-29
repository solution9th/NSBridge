package v1

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/solution9th/NSBridge/utils"
)

// GetTypes 获取可以添加的记录类型
func GetTypes(c *gin.Context) {

	lists := []string{
		"A", "AAAA", "MX", "CNAME", "TXT", "NS",
	}

	c.JSON(http.StatusOK, utils.ParseSuccessWithData(lists))
}

// GetLindeIDs 查看所有线路id列表
func GetLindeIDs(c *gin.Context) {

	c.String(http.StatusOK, "{\"err_code\":0,\"err_msg\":\"success\",\"data\":[%v,%v,%v,%v,%v,%v]}", lineContinental, lineISP, lineCountry, lineProvince, lineOutCity, lineChinaCityISP)
	return
}

// GetLineContinental 获得 洲 线路列表
func GetLineContinental(c *gin.Context) {
	c.String(http.StatusOK, "{\"err_code\":0,\"err_msg\":\"success\",\"data\":[%v]}", lineContinental)
	return
}

// GetLineISP 获得 ISP 线路列表
func GetLineISP(c *gin.Context) {
	c.String(http.StatusOK, "{\"err_code\":0,\"err_msg\":\"success\",\"data\":[%v]}", lineISP)
	return
}

// GetlineCountry 获得国家线路列表
func GetlineCountry(c *gin.Context) {
	c.String(http.StatusOK, "{\"err_code\":0,\"err_msg\":\"success\",\"data\":[%v]}", lineCountry)
	return
}

// GetLineProvince 获得省列表
func GetLineProvince(c *gin.Context) {
	c.String(http.StatusOK, "{\"err_code\":0,\"err_msg\":\"success\",\"data\":[%v]}", lineProvince)
	return
}

// GetLineOutCity 获得国外城市列表
func GetLineOutCity(c *gin.Context) {
	c.String(http.StatusOK, "{\"err_code\":0,\"err_msg\":\"success\",\"data\":[%v]}", lineOutCity)
	return
}

// GetLineChinaCityISP 获得国内城市及运营商列表
func GetLineChinaCityISP(c *gin.Context) {
	c.String(http.StatusOK, "{\"err_code\":0,\"err_msg\":\"success\",\"data\":[%v]}", lineChinaCityISP)
	return
}

const (
	// https://admin-dev.f1.zlibs.com/admin-api/dns-util/lines

	lineContinental = `{"id":0,"name":"default"},{"id":1,"name":"Africa"},{"id":2,"name":"Asia"},{"id":3,"name":"Europe"},{"id":4,"name":"North America"},{"id":5,"name":"South America"},{"id":6,"name":"Oceania"}`

	lineISP = `{"id":10,"name":"CHINAEDU"},{"id":11,"name":"ChinaMobile"},{"id":12,"name":"ChinaRailcom"},{"id":13,"name":"ChinaTelecom"},{"id":14,"name":"ChinaUnicom"},{"id":15,"name":"DRPENG"},{"id":16,"name":"FOUNDERBN"},{"id":17,"name":"TOPWAY"},{"id":18,"name":"WASU"}`

	lineCountry = `{"id":100,"name":"Afghanistan"},{"id":101,"name":"Algeria"},{"id":102,"name":"Angola"},{"id":103,"name":"Argentina"},{"id":104,"name":"Armenia"},{"id":105,"name":"Australia"},{"id":106,"name":"Austria"},{"id":107,"name":"Azerbaijan"},{"id":108,"name":"Bangladesh"},{"id":109,"name":"Belarus"},{"id":110,"name":"Belgium"},{"id":111,"name":"Bosnia and Herzegovina"},{"id":112,"name":"Brazil"},{"id":113,"name":"Bulgaria"},{"id":114,"name":"Cambodia"},{"id":115,"name":"Canada"},{"id":116,"name":"Chile"},{"id":117,"name":"China"},{"id":118,"name":"Hong Kong"},{"id":119,"name":"Taiwan"},{"id":120,"name":"Colombia"},{"id":121,"name":"Costa Rica"},{"id":122,"name":"Croatia"},{"id":123,"name":"Cyprus"},{"id":124,"name":"Czech"},{"id":125,"name":"Denmark"},{"id":126,"name":"Dominica"},{"id":127,"name":"Ecuador"},{"id":128,"name":"Egypt"},{"id":129,"name":"El Salvador"},{"id":130,"name":"Estonia"},{"id":131,"name":"Finland"},{"id":132,"name":"France"},{"id":133,"name":"Garner"},{"id":134,"name":"Georgia"},{"id":135,"name":"Germany"},{"id":136,"name":"Greece"},{"id":137,"name":"Guatemala"},{"id":138,"name":"Honduras"},{"id":139,"name":"Hungary"},{"id":140,"name":"Iceland"},{"id":141,"name":"India"},{"id":142,"name":"Indonesia"},{"id":143,"name":"Iran"},{"id":144,"name":"Iraq"},{"id":145,"name":"Ireland"},{"id":146,"name":"Israel"},{"id":147,"name":"Italy"},{"id":148,"name":"Japan"},{"id":149,"name":"Jordan"},{"id":150,"name":"Kazakhstan"},{"id":151,"name":"Kenya"},{"id":152,"name":"Kuwait"},{"id":153,"name":"Latvia"},{"id":154,"name":"Lebanon"},{"id":155,"name":"Lithuania"},{"id":156,"name":"Luxemburg"},{"id":157,"name":"Macedonia"},{"id":158,"name":"Malaysia"},{"id":159,"name":"Malta"},{"id":160,"name":"Mexico"},{"id":161,"name":"Moldova"},{"id":162,"name":"Morocco"},{"id":163,"name":"Nepal"},{"id":164,"name":"Netherlands"},{"id":165,"name":"New Zealand"},{"id":166,"name":"Nigeria"},{"id":167,"name":"Norway"},{"id":168,"name":"Pakistan"},{"id":169,"name":"Palestine"},{"id":170,"name":"Panama"},{"id":171,"name":"Peru"},{"id":172,"name":"Philippines"},{"id":173,"name":"Poland"},{"id":174,"name":"Portugal"},{"id":175,"name":"Puerto Rico"},{"id":176,"name":"Romania"},{"id":177,"name":"Russia"},{"id":178,"name":"Saudi Arabia"},{"id":179,"name":"Serbia"},{"id":180,"name":"Seychelles"},{"id":181,"name":"Singapore"},{"id":182,"name":"Slovakia"},{"id":183,"name":"Slovenia"},{"id":184,"name":"South Africa"},{"id":185,"name":"Spain"},{"id":186,"name":"Sweden"},{"id":187,"name":"Switzerland"},{"id":188,"name":"Syria"},{"id":189,"name":"Tanzania"},{"id":190,"name":"Thailand"},{"id":191,"name":"The Republic of Korea"},{"id":192,"name":"The United Arab Emirates"},{"id":193,"name":"Tunisia"},{"id":194,"name":"Turkey"},{"id":195,"name":"Ukraine"},{"id":196,"name":"United Kingdom"},{"id":197,"name":"United States"},{"id":198,"name":"Uruguay"},{"id":199,"name":"Venezuela"},{"id":200,"name":"Vietnam"}`

	lineProvince = `{"id":300,"name":"Anhui"},{"id":301,"name":"Beijing"},{"id":302,"name":"Chongqing"},{"id":303,"name":"Fujian"},{"id":304,"name":"Gansu"},{"id":305,"name":"Guangdong"},{"id":306,"name":"Guangxi"},{"id":307,"name":"Guizhou"},{"id":308,"name":"Hainan"},{"id":309,"name":"Hebei"},{"id":310,"name":"Heilongjiang"},{"id":311,"name":"Henan"},{"id":312,"name":"Hubei"},{"id":313,"name":"Hunan"},{"id":314,"name":"Jiangsu"},{"id":315,"name":"Jiangxi"},{"id":316,"name":"Jilin"},{"id":317,"name":"Liaoning"},{"id":318,"name":"Nei Mongol"},{"id":319,"name":"Ningxia"},{"id":320,"name":"Qinghai"},{"id":321,"name":"Shaanxi"},{"id":322,"name":"Shandong"},{"id":323,"name":"Shanghai"},{"id":324,"name":"Shanxi"},{"id":325,"name":"Sichuan"},{"id":326,"name":"Tianjin"},{"id":327,"name":"Xinjiang"},{"id":328,"name":"Xizang"},{"id":329,"name":"Yunnan"},{"id":330,"name":"Zhejiang"}`

	lineOutCity = `{"id":400,"name":"Alabama"},{"id":401,"name":"Alaska"},{"id":402,"name":"Arizona"},{"id":403,"name":"Arkansas"},{"id":404,"name":"California"},{"id":405,"name":"Colorado"},{"id":406,"name":"Connecticut"},{"id":407,"name":"Delaware"},{"id":408,"name":"District of Columbia"},{"id":409,"name":"Florida"},{"id":410,"name":"Georgia"},{"id":411,"name":"Hawaii"},{"id":412,"name":"Idaho"},{"id":413,"name":"Illinois"},{"id":414,"name":"Indiana"},{"id":415,"name":"Iowa"},{"id":416,"name":"Kansas"},{"id":417,"name":"Kentucky"},{"id":418,"name":"Louisiana"},{"id":419,"name":"Maine"},{"id":420,"name":"Maryland"},{"id":421,"name":"Massachusetts"},{"id":422,"name":"Michigan"},{"id":423,"name":"Minnesota"},{"id":424,"name":"Mississippi"},{"id":425,"name":"Missouri"},{"id":426,"name":"Montana"},{"id":427,"name":"Nebraska"},{"id":428,"name":"Nevada"},{"id":429,"name":"New Hampshire"},{"id":430,"name":"New Jersey"},{"id":431,"name":"New Mexico"},{"id":432,"name":"New York"},{"id":433,"name":"North Carolina"},{"id":434,"name":"North Dakota"},{"id":435,"name":"Ohio"},{"id":436,"name":"Oklahoma"},{"id":437,"name":"Oregon"},{"id":438,"name":"Pennsylvania"},{"id":439,"name":"Rhode Island"},{"id":440,"name":"South Carolina"},{"id":441,"name":"South Dakota"},{"id":442,"name":"Tennessee"},{"id":443,"name":"Texas"},{"id":444,"name":"Utah"},{"id":445,"name":"Vermont"},{"id":446,"name":"Virginia"},{"id":447,"name":"Washington"},{"id":448,"name":"West Virginia"},{"id":449,"name":"Wisconsin"},{"id":450,"name":"Wyoming"}`

	lineChinaCityISP = `{"id":1000,"name":"Anhui CHINAEDU"},{"id":1001,"name":"Anhui ChinaMobile"},{"id":1002,"name":"Anhui ChinaRailcom"},{"id":1003,"name":"Anhui ChinaTelecom"},{"id":1004,"name":"Anhui ChinaUnicom"},{"id":1005,"name":"Anhui DRPENG"},{"id":1006,"name":"Beijing CHINAEDU"},{"id":1007,"name":"Beijing ChinaMobile"},{"id":1008,"name":"Beijing ChinaRailcom"},{"id":1009,"name":"Beijing ChinaTelecom"},{"id":1010,"name":"Beijing ChinaUnicom"},{"id":1011,"name":"Beijing DRPENG"},{"id":1012,"name":"China CHINAEDU"},{"id":1013,"name":"China ChinaMobile"},{"id":1014,"name":"China ChinaRailcom"},{"id":1015,"name":"China ChinaTelecom"},{"id":1016,"name":"China ChinaUnicom"},{"id":1017,"name":"China DRPENG"},{"id":1018,"name":"Chongqing CHINAEDU"},{"id":1019,"name":"Chongqing ChinaMobile"},{"id":1020,"name":"Chongqing ChinaRailcom"},{"id":1021,"name":"Chongqing ChinaTelecom"},{"id":1022,"name":"Chongqing ChinaUnicom"},{"id":1023,"name":"Chongqing DRPENG"},{"id":1024,"name":"Fujian CHINAEDU"},{"id":1025,"name":"Fujian ChinaMobile"},{"id":1026,"name":"Fujian ChinaRailcom"},{"id":1027,"name":"Fujian ChinaTelecom"},{"id":1028,"name":"Fujian ChinaUnicom"},{"id":1029,"name":"Fujian DRPENG"},{"id":1030,"name":"Gansu CHINAEDU"},{"id":1031,"name":"Gansu ChinaMobile"},{"id":1032,"name":"Gansu ChinaRailcom"},{"id":1033,"name":"Gansu ChinaTelecom"},{"id":1034,"name":"Gansu ChinaUnicom"},{"id":1035,"name":"Gansu DRPENG"},{"id":1036,"name":"Guangdong CHINAEDU"},{"id":1037,"name":"Guangdong ChinaMobile"},{"id":1038,"name":"Guangdong ChinaRailcom"},{"id":1039,"name":"Guangdong ChinaTelecom"},{"id":1040,"name":"Guangdong ChinaUnicom"},{"id":1041,"name":"Guangdong DRPENG"},{"id":1042,"name":"Guangxi CHINAEDU"},{"id":1043,"name":"Guangxi ChinaMobile"},{"id":1044,"name":"Guangxi ChinaRailcom"},{"id":1045,"name":"Guangxi ChinaTelecom"},{"id":1046,"name":"Guangxi ChinaUnicom"},{"id":1047,"name":"Guangxi DRPENG"},{"id":1048,"name":"Guizhou CHINAEDU"},{"id":1049,"name":"Guizhou ChinaMobile"},{"id":1050,"name":"Guizhou ChinaRailcom"},{"id":1051,"name":"Guizhou ChinaTelecom"},{"id":1052,"name":"Guizhou ChinaUnicom"},{"id":1053,"name":"Guizhou DRPENG"},{"id":1054,"name":"Hainan CHINAEDU"},{"id":1055,"name":"Hainan ChinaMobile"},{"id":1056,"name":"Hainan ChinaRailcom"},{"id":1057,"name":"Hainan ChinaTelecom"},{"id":1058,"name":"Hainan ChinaUnicom"},{"id":1059,"name":"Hainan DRPENG"},{"id":1060,"name":"Hebei CHINAEDU"},{"id":1061,"name":"Hebei ChinaMobile"},{"id":1062,"name":"Hebei ChinaRailcom"},{"id":1063,"name":"Hebei ChinaTelecom"},{"id":1064,"name":"Hebei ChinaUnicom"},{"id":1065,"name":"Hebei DRPENG"},{"id":1066,"name":"Heilongjiang CHINAEDU"},{"id":1067,"name":"Heilongjiang ChinaMobile"},{"id":1068,"name":"Heilongjiang ChinaRailcom"},{"id":1069,"name":"Heilongjiang ChinaTelecom"},{"id":1070,"name":"Heilongjiang ChinaUnicom"},{"id":1071,"name":"Heilongjiang DRPENG"},{"id":1072,"name":"Henan CHINAEDU"},{"id":1073,"name":"Henan ChinaMobile"},{"id":1074,"name":"Henan ChinaRailcom"},{"id":1075,"name":"Henan ChinaTelecom"},{"id":1076,"name":"Henan ChinaUnicom"},{"id":1077,"name":"Henan DRPENG"},{"id":1078,"name":"Hubei CHINAEDU"},{"id":1079,"name":"Hubei ChinaMobile"},{"id":1080,"name":"Hubei ChinaRailcom"},{"id":1081,"name":"Hubei ChinaTelecom"},{"id":1082,"name":"Hubei ChinaUnicom"},{"id":1083,"name":"Hubei DRPENG"},{"id":1084,"name":"Hunan CHINAEDU"},{"id":1085,"name":"Hunan ChinaMobile"},{"id":1086,"name":"Hunan ChinaRailcom"},{"id":1087,"name":"Hunan ChinaTelecom"},{"id":1088,"name":"Hunan ChinaUnicom"},{"id":1089,"name":"Hunan DRPENG"},{"id":1090,"name":"Jiangsu CHINAEDU"},{"id":1091,"name":"Jiangsu ChinaMobile"},{"id":1092,"name":"Jiangsu ChinaRailcom"},{"id":1093,"name":"Jiangsu ChinaTelecom"},{"id":1094,"name":"Jiangsu ChinaUnicom"},{"id":1095,"name":"Jiangsu DRPENG"},{"id":1096,"name":"Jiangxi CHINAEDU"},{"id":1097,"name":"Jiangxi ChinaMobile"},{"id":1098,"name":"Jiangxi ChinaRailcom"},{"id":1099,"name":"Jiangxi ChinaTelecom"},{"id":1100,"name":"Jiangxi ChinaUnicom"},{"id":1101,"name":"Jiangxi DRPENG"},{"id":1102,"name":"Jilin CHINAEDU"},{"id":1103,"name":"Jilin ChinaMobile"},{"id":1104,"name":"Jilin ChinaRailcom"},{"id":1105,"name":"Jilin ChinaTelecom"},{"id":1106,"name":"Jilin ChinaUnicom"},{"id":1107,"name":"Jilin DRPENG"},{"id":1108,"name":"Liaoning CHINAEDU"},{"id":1109,"name":"Liaoning ChinaMobile"},{"id":1110,"name":"Liaoning ChinaRailcom"},{"id":1111,"name":"Liaoning ChinaTelecom"},{"id":1112,"name":"Liaoning ChinaUnicom"},{"id":1113,"name":"Liaoning DRPENG"},{"id":1114,"name":"Nei Mongol CHINAEDU"},{"id":1115,"name":"Nei Mongol ChinaMobile"},{"id":1116,"name":"Nei Mongol ChinaRailcom"},{"id":1117,"name":"Nei Mongol ChinaTelecom"},{"id":1118,"name":"Nei Mongol ChinaUnicom"},{"id":1119,"name":"Nei Mongol DRPENG"},{"id":1120,"name":"Ningxia CHINAEDU"},{"id":1121,"name":"Ningxia ChinaMobile"},{"id":1122,"name":"Ningxia ChinaRailcom"},{"id":1123,"name":"Ningxia ChinaTelecom"},{"id":1124,"name":"Ningxia ChinaUnicom"},{"id":1125,"name":"Qinghai CHINAEDU"},{"id":1126,"name":"Qinghai ChinaMobile"},{"id":1127,"name":"Qinghai ChinaRailcom"},{"id":1128,"name":"Qinghai ChinaTelecom"},{"id":1129,"name":"Qinghai ChinaUnicom"},{"id":1130,"name":"Shaanxi CHINAEDU"},{"id":1131,"name":"Shaanxi ChinaMobile"},{"id":1132,"name":"Shaanxi ChinaRailcom"},{"id":1133,"name":"Shaanxi ChinaTelecom"},{"id":1134,"name":"Shaanxi ChinaUnicom"},{"id":1135,"name":"Shaanxi DRPENG"},{"id":1136,"name":"Shandong CHINAEDU"},{"id":1137,"name":"Shandong ChinaMobile"},{"id":1138,"name":"Shandong ChinaRailcom"},{"id":1139,"name":"Shandong ChinaTelecom"},{"id":1140,"name":"Shandong ChinaUnicom"},{"id":1141,"name":"Shandong DRPENG"},{"id":1142,"name":"Shanghai CHINAEDU"},{"id":1143,"name":"Shanghai ChinaMobile"},{"id":1144,"name":"Shanghai ChinaRailcom"},{"id":1145,"name":"Shanghai ChinaTelecom"},{"id":1146,"name":"Shanghai ChinaUnicom"},{"id":1147,"name":"Shanghai DRPENG"},{"id":1148,"name":"Shanxi CHINAEDU"},{"id":1149,"name":"Shanxi ChinaMobile"},{"id":1150,"name":"Shanxi ChinaRailcom"},{"id":1151,"name":"Shanxi ChinaTelecom"},{"id":1152,"name":"Shanxi ChinaUnicom"},{"id":1153,"name":"Shanxi DRPENG"},{"id":1154,"name":"Sichuan CHINAEDU"},{"id":1155,"name":"Sichuan ChinaMobile"},{"id":1156,"name":"Sichuan ChinaRailcom"},{"id":1157,"name":"Sichuan ChinaTelecom"},{"id":1158,"name":"Sichuan ChinaUnicom"},{"id":1159,"name":"Sichuan DRPENG"},{"id":1160,"name":"Tianjin CHINAEDU"},{"id":1161,"name":"Tianjin ChinaMobile"},{"id":1162,"name":"Tianjin ChinaRailcom"},{"id":1163,"name":"Tianjin ChinaTelecom"},{"id":1164,"name":"Tianjin ChinaUnicom"},{"id":1165,"name":"Tianjin DRPENG"},{"id":1166,"name":"Xinjiang CHINAEDU"},{"id":1167,"name":"Xinjiang ChinaMobile"},{"id":1168,"name":"Xinjiang ChinaRailcom"},{"id":1169,"name":"Xinjiang ChinaTelecom"},{"id":1170,"name":"Xinjiang ChinaUnicom"},{"id":1171,"name":"Xizang CHINAEDU"},{"id":1172,"name":"Xizang ChinaMobile"},{"id":1173,"name":"Xizang ChinaRailcom"},{"id":1174,"name":"Xizang ChinaTelecom"},{"id":1175,"name":"Xizang ChinaUnicom"},{"id":1176,"name":"Yunnan CHINAEDU"},{"id":1177,"name":"Yunnan ChinaMobile"},{"id":1178,"name":"Yunnan ChinaRailcom"},{"id":1179,"name":"Yunnan ChinaTelecom"},{"id":1180,"name":"Yunnan ChinaUnicom"},{"id":1181,"name":"Yunnan DRPENG"},{"id":1182,"name":"Zhejiang CHINAEDU"},{"id":1183,"name":"Zhejiang ChinaMobile"},{"id":1184,"name":"Zhejiang ChinaRailcom"},{"id":1185,"name":"Zhejiang ChinaTelecom"},{"id":1186,"name":"Zhejiang ChinaUnicom"},{"id":1187,"name":"Zhejiang DRPENG"}`
)
