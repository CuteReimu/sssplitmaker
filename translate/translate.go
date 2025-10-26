package translate

import (
	_ "embed"
	"encoding/json"
	"slices"
	"strings"
)

//go:embed splits.json
var splitsJson []byte

var SplitsHtml string

func init() {
	var splitCache []*SplitData
	err := json.Unmarshal(splitsJson, &splitCache)
	if err != nil {
		panic(err)
	}
	for _, split := range splitCache {
		split.Description = GetSplitDescriptionByID(split.ID)
	}
	buf, err := json.Marshal(splitCache)
	if err != nil {
		panic(err)
	}
	SplitsHtml = string(buf)
}

type SplitData struct {
	ID          string `json:"key"`
	Description string `json:"translate"`
	Desc        string `json:"description"`
	Tooltip     string `json:"tooltip"`
	Alias       any    `json:"alias,omitempty"`
}

func GetSplitDescriptionByID(id string) string {
	if v, ok := cacheAliases[id]; ok {
		id = v
	}
	for _, split := range SplitsCache {
		if split.ID == id {
			return split.Description
		}
	}
	return ""
}

func GetIndexByID(id string) int {
	if v, ok := cacheAliases[id]; ok {
		id = v
	}
	for i, split := range SplitsCache {
		if split.ID == id {
			return i
		}
	}
	return -1
}

func GetIDByDescription(desc string) string {
	for _, split := range SplitsCache {
		if split.Description == desc {
			return split.ID
		}
	}
	return ""
}

func init() {
	typeOrder := map[string]int{
		"（其它）": -100,
		"（开始）": -90,
		"（结束）": -80,
		"（旋律）": -70,
		"（道具）": -60,
	}
	slices.SortStableFunc(SplitsCache, func(a, b *SplitData) int {
		i := strings.LastIndex(a.Description, "（")
		j := strings.LastIndex(b.Description, "（")
		if i == -1 {
			i = len(a.Description)
		}
		if j == -1 {
			j = len(b.Description)
		}
		ad, bd := a.Description[i:], b.Description[j:]
		aScore, bScore := typeOrder[ad], typeOrder[bd]
		if aScore != bScore {
			return aScore - bScore
		}
		d := strings.Compare(ad, bd)
		if d != 0 {
			return d
		}
		return strings.Compare(a.Description, b.Description)
	})
}

var SplitsCache = []*SplitData{
	{ID: "ManualSplit", Description: "手动分割（其它）"},
	{ID: "StartNewGame", Description: "开始新游戏（开始）"},
	{ID: "EndingSplit", Description: "任意结束（结束）"},
	{ID: "EndingA", Description: "苍白之母结局（结束）"},
	{ID: "Menu", Description: "主菜单（菜单）"},
	{ID: "PlayerDeath", Description: "死亡（事件）"},
	{ID: "AnyTransition", Description: "任意切图（切图）"},
	{ID: "MossMother", Description: "苔藓之母（Boss）"},
	{ID: "MossMotherTrans", Description: "苔藓之母（切图）"},
	{ID: "SilkSpear", Description: "丝之矛（技能）"},
	{ID: "SilkSpearTrans", Description: "丝之矛（切图）"},
	{ID: "BellBeast", Description: "钟道兽（Boss）"},
	{ID: "BellBeastTrans", Description: "钟道兽（切图）"},
	{ID: "MarrowBell", Description: "敲钟-髓骨洞窟（事件）"},
	{ID: "SwiftStep", Description: "冲刺（技能）"},
	{ID: "SwiftStepTrans", Description: "冲刺（切图）"},
	{ID: "Lace1", Description: "蕾丝（Boss）"},
	{ID: "Lace1Trans", Description: "蕾丝（切图）"},
	{ID: "DeepDocksBell", Description: "敲钟-深坞（事件）"},
	{ID: "DriftersCloak", Description: "斗篷（技能）"},
	{ID: "DriftersCloakTrans", Description: "斗篷（切图）"},
	{ID: "FourthChorus", Description: "第四圣咏团（Boss）"},
	{ID: "EnterGreymoor", Description: "进入灰沼（切图）"},
	{ID: "GreymoorBell", Description: "敲钟-灰沼（事件）"},
	{ID: "Moorwing", Description: "荒沼翼主（Boss）"},
	{ID: "MoorwingTrans", Description: "荒沼翼主（切图）"},
	{ID: "ClingGrip", Description: "爬墙（技能）"},
	{ID: "ClingGripTrans", Description: "爬墙（切图）"},
	{ID: "ShellwoodBell", Description: "敲钟-甲木林（事件）"},
	{ID: "Widow", Description: "黑寡妇（Boss）"},
	{ID: "BellhartBell", Description: "敲钟-钟心镇（事件）"},
	{ID: "LastJudge", Description: "末代裁决者（Boss）"},
	{ID: "EnterMist", Description: "进入迷雾（切图）"},
	{ID: "LeaveMist", Description: "离开迷雾（切图）"},
	{ID: "Phantom", Description: "幽影（Boss）"},
	{ID: "Act2Started", Description: "第二幕开始（事件）"},
	{ID: "CogworkDancers", Description: "机枢舞者（Boss）"},
	{ID: "WhisperingVaultsArena", Description: "低语书库遭遇战（小Boss）"},
	{ID: "Trobbio", Description: "特罗比奥（Boss）"},
	{ID: "TrobbioTrans", Description: "特罗比奥（切图）"},
	{ID: "Clawline", Description: "飞针冲刺（技能）"},
	{ID: "EnterHighHalls", Description: "进入高庭（切图）"},
	{ID: "EnterHighHallsArena", Description: "进入高庭遭遇战（切图）"},
	{ID: "HighHallsArena", Description: "高庭遭遇战（小Boss）"},
	{ID: "Lace2", Description: "蕾丝2（Boss）"},
	{ID: "VaultkeepersMelody", Description: "管理员旋律（旋律）"},
	{ID: "VaultkeepersMelodyTrans", Description: "管理员旋律（切图）"},
	{ID: "ArchitectsMelody", Description: "建筑师旋律（旋律）"},
	{ID: "ArchitectsMelodyTrans", Description: "建筑师旋律（切图）"},
	{ID: "ConductorsMelody", Description: "指挥家旋律（旋律）"},
	{ID: "ConductorsMelodyTrans", Description: "指挥家旋律（切图）"},
	{ID: "UnlockedMelodyLift", Description: "解锁三重旋律电梯（事件）"},
	{ID: "NeedleUpgrade1", Description: "织针升级1（升级）"},
	{ID: "NeedleUpgrade2", Description: "织针升级2（升级）"},
	{ID: "NeedleUpgrade3", Description: "织针升级3（升级）"},
	{ID: "NeedleUpgrade4", Description: "织针升级4（升级）"},
	{ID: "SavedFleaHuntersMarch", Description: "救跳蚤-猎者小径（跳蚤）"},
	{ID: "SavedFleaBellhart", Description: "救跳蚤-钟心镇（跳蚤）"},
	{ID: "SavedFleaMarrow", Description: "救跳蚤-髓骨洞窟（跳蚤）"},
	{ID: "SavedFleaDeepDocksSprint", Description: "救跳蚤-深坞-冲刺（跳蚤）"},
	{ID: "SavedFleaFarFieldsPilgrimsRest", Description: "救跳蚤-远野-朝圣者憩所（跳蚤）"},
	{ID: "SavedFleaFarFieldsTrap", Description: "救跳蚤-远野-陷阱（跳蚤）"},
	{ID: "SavedFleaSandsOfKarak", Description: "救跳蚤-卡拉卡沙川（跳蚤）"},
	{ID: "SavedFleaBlastedSteps", Description: "救跳蚤-卡拉卡沙川-蚀阶（跳蚤）"},
	{ID: "SavedFleaWormways", Description: "救跳蚤-沙噬虫道（跳蚤）"},
	{ID: "SavedFleaDeepDocksArena", Description: "救跳蚤-深坞-遭遇战（跳蚤）"},
	{ID: "SavedFleaDeepDocksBellway", Description: "救跳蚤-深坞-钟道（跳蚤）"},
	{ID: "SavedFleaBilewaterOrgan", Description: "救跳蚤-废鸣管风琴（跳蚤）"},
	{ID: "SavedFleaSinnersRoad", Description: "救跳蚤-罪途（跳蚤）"},
	{ID: "SavedFleaGreymoorRoof", Description: "救跳蚤-灰沼-屋顶（跳蚤）"},
	{ID: "SavedFleaGreymoorLake", Description: "救跳蚤-灰沼-湖（跳蚤）"},
	{ID: "SavedFleaWhisperingVaults", Description: "救跳蚤-低语书库（跳蚤）"},
	{ID: "SavedFleaSongclave", Description: "救跳蚤-低语书库-圣歌盟地（跳蚤）"},
	{ID: "SavedFleaMountFay", Description: "救跳蚤-费耶山（跳蚤）"},
	{ID: "SavedFleaBilewaterTrap", Description: "救跳蚤-腐汁泽-陷阱（跳蚤）"},
	{ID: "SavedFleaBilewaterThieves", Description: "救跳蚤-腐汁泽-小偷（跳蚤）"},
	{ID: "SavedFleaShellwood", Description: "救跳蚤-甲木林（跳蚤）"},
	{ID: "SavedFleaSlabBellway", Description: "救跳蚤-罪石监狱-钟道（跳蚤）"},
	{ID: "SavedFleaSlabCage", Description: "救跳蚤-罪石监狱-笼子（跳蚤）"},
	{ID: "SavedFleaChoralChambersWind", Description: "救跳蚤-圣咏殿-风（跳蚤）"},
	{ID: "SavedFleaChoralChambersCage", Description: "救跳蚤-圣咏殿-笼子（跳蚤）"},
	{ID: "SavedFleaUnderworksCauldron", Description: "救跳蚤-圣堡工厂-大熔釜（跳蚤）"},
	{ID: "SavedFleaUnderworksWispThicket", Description: "救跳蚤-圣堡工厂-火灵竹丛（跳蚤）"},
	{ID: "SavedFleaGiantFlea", Description: "击败大跳蚤（跳蚤）"},
	{ID: "SavedFleaVog", Description: "沃格（跳蚤）"},
	{ID: "SavedFleaKratt", Description: "救跳蚤-灰沼-克拉特（跳蚤）"},
	{ID: "PutrifiedDuctsStation", Description: "钟道-腐殖渠（钟道）"},
	{ID: "BellhartStation", Description: "钟道-钟心镇（钟道）"},
	{ID: "FarFieldsStation", Description: "钟道-远野（钟道）"},
	{ID: "GrandBellwayStation", Description: "钟道-圣堡钟道（钟道）"},
	{ID: "BlastedStepsStation", Description: "钟道-蚀阶（钟道）"},
	{ID: "DeepDocksStation", Description: "钟道-深坞（钟道）"},
	{ID: "GreymoorStation", Description: "钟道-灰沼（钟道）"},
	{ID: "SlabStation", Description: "钟道-罪石牢狱（钟道）"},
	{ID: "BilewaterStation", Description: "钟道-腐汁泽（钟道）"},
	{ID: "ShellwoodStation", Description: "钟道-甲木林（钟道）"},
	{ID: "ChoralChambersTube", Description: "圣脉枢管-圣咏殿（圣脉枢管）"},
	{ID: "UnderworksTube", Description: "圣脉枢管-圣堡工厂（圣脉枢管）"},
	{ID: "GrandBellwayTube", Description: "圣脉枢管-圣堡钟道（圣脉枢管）"},
	{ID: "HighHallsTube", Description: "圣脉枢管-高庭（圣脉枢管）"},
	{ID: "SongclaveTube", Description: "圣脉枢管-始源钟殿（圣脉枢管）"},
	{ID: "MemoriumTube", Description: "圣脉枢管-忆廊（圣脉枢管）"},
	{ID: "SeenShakraBonebottom", Description: "制图师-骸底镇（NPC）"},
	{ID: "SeenShakraMarrow", Description: "制图师-髓骨洞窟（NPC）"},
	{ID: "SeenShakraDeepDocks", Description: "制图师-深坞（NPC）"},
	{ID: "SeenShakraFarFields", Description: "制图师-远野（NPC）"},
	{ID: "SeenShakraWormways", Description: "制图师-沙噬虫道（NPC）"},
	{ID: "SeenShakraGreymoor", Description: "制图师-灰沼（NPC）"},
	{ID: "SeenShakraBellhart", Description: "制图师-钟心镇（NPC）"},
	{ID: "SeenShakraShellwood", Description: "制图师-甲木林（NPC）"},
	{ID: "SeenShakraHuntersMarch", Description: "制图师-猎者小径（NPC）"},
	{ID: "SeenShakraBlastedSteps", Description: "制图师-蚀阶（NPC）"},
	{ID: "SeenShakraSinnersRoad", Description: "制图师-罪途（NPC）"},
	{ID: "SeenShakraMountFay", Description: "制图师-费耶山（NPC）"},
	{ID: "SeenShakraBilewater", Description: "制图师-腐汁泽（NPC）"},
	{ID: "SeenShakraSandsOfKarak", Description: "制图师-卡拉卡沙川（NPC）"},
	{ID: "MetJubilanaEnclave", Description: "朱比拉娜-圣歌盟地（NPC）"},
	{ID: "MetShermaEnclave", Description: "谢尔玛-圣歌盟地（NPC）"},
	{ID: "UnlockedPrinceCage", Description: "绿王子-罪途（事件）"},
	{ID: "GreenPrinceInVerdania", Description: "绿王子-机枢核心（事件）"},
	{ID: "SeenFleatopiaEmpty", Description: "蚤托邦（事件）"},
	{ID: "FaydownCloak", Description: "二段跳（技能）"},
	{ID: "SilkSoar", Description: "灵丝升腾（技能）"},
	{ID: "RedMemory", Description: "赤红忆境（事件）"},
	{ID: "BellhouseKeyConversation", Description: "钟居钥匙（NPC）"},
	{ID: "Act1Start", Description: "第一幕开始（开始）"},
	{ID: "EnterWormways", Description: "进入沙噬虫道（切图）"},
	{ID: "EnterFarFields", Description: "进入远野（切图）"},
	{ID: "EnterShellwood", Description: "进入甲木林（切图）"},
	{ID: "EnterBellhart", Description: "进入钟心镇（切图）"},
	{ID: "ReaperCrestTrans", Description: "收割者纹章（切图）"},
	{ID: "HeartNyleth", Description: "尼莱斯之心（道具）"},
	{ID: "HeartKhann", Description: "卡汗之心（道具）"},
	{ID: "HeartKarmelita", Description: "卡梅莉塔之心（道具）"},
	{ID: "HeartClover", Description: "双生之心（道具）"},
	{ID: "VerdaniaOrbsCollected", Description: "翠庭球收集完成（事件）"},
	{ID: "Forebrothers", Description: "监工兄弟（Boss）"},
	{ID: "Groal", Description: "格洛（Boss）"},
	{ID: "GreatConchflies", Description: "巨型螺蝇（Boss）"},
	{ID: "GreatConchfliesTrans", Description: "巨型螺蝇（切图）"},
	{ID: "SavageBeastfly1", Description: "残暴的兽蝇1（Boss）"},
	{ID: "SavageBeastfly2", Description: "暴怒的兽蝇2（Boss）"},
	{ID: "CaravanTroupeGreymoor", Description: "跳蚤剧团移动-灰沼（事件）"},
	{ID: "CaravanTroupeFleatopia", Description: "跳蚤剧团移动-蚤托邦（事件）"},
	{ID: "SoldRelic", Description: "卖古董（事件）"},
	{ID: "CollectedWhitewardKey", Description: "白愈钥匙（道具）"},
	{ID: "PavoTimePassed", Description: "与帕沃对话后休息（坐椅子）"},
	{ID: "SongclaveBell", Description: "敲钟-圣歌盟地（事件）"},
	{ID: "Voltvyrm", Description: "伏特维姆（Boss）"},
	{ID: "SkullTyrant1", Description: "骷髅暴君1（Boss）"},
	{ID: "ShermaReturned", Description: "谢尔玛归来（NPC）"},
	{ID: "JubilanaRescuedMemorium", Description: "救朱比拉娜-忆廊（事件）"},
	{ID: "JubilanaRescuedChoralChambers", Description: "救朱比拉娜-圣咏殿（事件）"},
	{ID: "SilkAndSoulOffered", Description: "灵丝与灵魂前置达成（事件）"},
	{ID: "SoulSnareReady", Description: "灵丝陷阱就绪（事件）"},
	{ID: "Seth", Description: "赛斯（Boss）"},
	{ID: "AbyssEscape", Description: "深渊逃脱完成（事件）"},
	{ID: "BallowMoved", Description: "巴洛前往潜钟（事件）"},
	{ID: "Act3Started", Description: "第三幕开始（事件）"},
	{ID: "EnterMosshome", Description: "进入苔栖乡（切图）"},
	{ID: "BoneBottomSimpleKey", Description: "骸底镇简单钥匙（道具）"},
	{ID: "EnterUpperWormways", Description: "进入沙噬虫道上层（切图）"},
	{ID: "Sharpdart", Description: "丝刃镖（技能）"},
	{ID: "SharpdartTrans", Description: "丝刃镖（切图）"},
	{ID: "EnterHuntersMarch", Description: "进入猎者小径（切图）"},
	{ID: "HuntersMarchPostMiddleArenaTransition", Description: "进入猎者小径-中间遭遇战（切图）"},
	{ID: "ThreadStorm", Description: "灵丝风暴（技能）"},
	{ID: "ThreadStormTrans", Description: "灵丝风暴（切图）"},
	{ID: "EnterBlastedSteps", Description: "进入蚀阶（切图）"},
	{ID: "EnterLastJudge", Description: "进入末日裁决者（切图）"},
	{ID: "EnterSinnersRoad", Description: "进入罪途（切图）"},
	{ID: "EnterBilewater", Description: "进入腐汁泽（切图）"},
	{ID: "EnterExhaustOrgan", Description: "进入废鸣管风琴（切图）"},
	{ID: "CrossStitch", Description: "十字缝（技能）"},
	{ID: "CrossStitchTrans", Description: "十字缝（切图）"},
	{ID: "RuneRage", Description: "符文之怒（技能）"},
	{ID: "RuneRageTrans", Description: "符文之怒（切图）"},
	{ID: "EnterMemorium", Description: "进入忆廊（切图）"},
	{ID: "PaleNails", Description: "苍白之爪（技能）"},
	{ID: "PaleNailsTrans", Description: "苍白之爪（切图）"},
	{ID: "MaskShard1", Description: "面具碎片1（碎片）"},
	{ID: "MaskShard2", Description: "面具碎片2（碎片）"},
	{ID: "MaskShard3", Description: "面具碎片3（碎片）"},
	{ID: "Mask1", Description: "面具碎片4（碎片）"},
	{ID: "MaskShard5", Description: "面具碎片5（碎片）"},
	{ID: "MaskShard6", Description: "面具碎片6（碎片）"},
	{ID: "MaskShard7", Description: "面具碎片7（碎片）"},
	{ID: "Mask2", Description: "面具碎片8（碎片）"},
	{ID: "MaskShard9", Description: "面具碎片9（碎片）"},
	{ID: "MaskShard10", Description: "面具碎片10（碎片）"},
	{ID: "MaskShard11", Description: "面具碎片11（碎片）"},
	{ID: "Mask3", Description: "面具碎片12（碎片）"},
	{ID: "MaskShard13", Description: "面具碎片13（碎片）"},
	{ID: "MaskShard14", Description: "面具碎片14（碎片）"},
	{ID: "MaskShard15", Description: "面具碎片15（碎片）"},
	{ID: "Mask4", Description: "面具碎片16（碎片）"},
	{ID: "MaskShard17", Description: "面具碎片17（碎片）"},
	{ID: "MaskShard18", Description: "面具碎片18（碎片）"},
	{ID: "MaskShard19", Description: "面具碎片19（碎片）"},
	{ID: "Mask5", Description: "面具碎片20（碎片）"},
	{ID: "SpoolFragment1", Description: "灵丝轴碎片1（碎片）"},
	{ID: "Spool1", Description: "灵丝轴碎片2（碎片）"},
	{ID: "SpoolFragment3", Description: "灵丝轴碎片3（碎片）"},
	{ID: "Spool2", Description: "灵丝轴碎片4（碎片）"},
	{ID: "SpoolFragment5", Description: "灵丝轴碎片5（碎片）"},
	{ID: "Spool3", Description: "灵丝轴碎片6（碎片）"},
	{ID: "SpoolFragment7", Description: "灵丝轴碎片7（碎片）"},
	{ID: "Spool4", Description: "灵丝轴碎片8（碎片）"},
	{ID: "SpoolFragment9", Description: "灵丝轴碎片9（碎片）"},
	{ID: "Spool5", Description: "灵丝轴碎片10（碎片）"},
	{ID: "SpoolFragment11", Description: "灵丝轴碎片11（碎片）"},
	{ID: "Spool6", Description: "灵丝轴碎片12（碎片）"},
	{ID: "SpoolFragment13", Description: "灵丝轴碎片13（碎片）"},
	{ID: "Spool7", Description: "灵丝轴碎片14（碎片）"},
	{ID: "SpoolFragment15", Description: "灵丝轴碎片15（碎片）"},
	{ID: "Spool8", Description: "灵丝轴碎片16（碎片）"},
	{ID: "SpoolFragment17", Description: "灵丝轴碎片17（碎片）"},
	{ID: "Spool9", Description: "灵丝轴碎片18（碎片）"},
	{ID: "ReaperCrest", Description: "收割者纹章（纹章）"},
	{ID: "WandererCrest", Description: "漫游者纹章（纹章）"},
	{ID: "WandererCrestTrans", Description: "漫游者纹章（切图）"},
	{ID: "BeastCrest", Description: "野兽纹章（纹章）"},
	{ID: "BeastCrestTrans", Description: "野兽纹章（切图）"},
	{ID: "ArchitectCrest", Description: "建筑师纹章（纹章）"},
	{ID: "ArchitectCrestTrans", Description: "建筑师纹章（切图）"},
	{ID: "ShamanCrest", Description: "萨满纹章（纹章）"},
	{ID: "ShamanCrestTrans", Description: "萨满纹章（切图）"},
	{ID: "EnterBellEater", Description: "进入噬钟者（切图）"},
	{ID: "BeastlingCall", Description: "唤兽曲（技能）"},
	{ID: "EnterNylethMemory", Description: "进入尼莱斯记忆（切图）"},
	{ID: "EnterKhannMemory", Description: "进入卡汗记忆（切图）"},
	{ID: "EnterKarmelitaMemory", Description: "进入卡梅莉塔记忆（切图）"},
	{ID: "EnterVerdaniaMemory", Description: "进入翠庭记忆（切图）"},
	{ID: "EnterVerdaniaCastle", Description: "进入翠庭城堡（切图）"},
	{ID: "VerdaniaLakeFountainOrbs", Description: "翠庭-喷泉（事件）"},
	{ID: "EnterSeth", Description: "进入赛斯战斗（切图）"},
	{ID: "TransitionExcludingDiscontinuities", Description: "任意切图不含非连续（切图）"},
	{ID: "EnterBoneBottom", Description: "进入骸底镇（切图）"},
	{ID: "NeedleStrike", Description: "蓄力斩（技能）"},
	{ID: "NeedleStrikeTrans", Description: "蓄力斩（切图）"},
	{ID: "EnterCitadelFrontGate", Description: "进入圣堡正门（切图）"},
	{ID: "EnterWhisperingVaults", Description: "进入低语书库（切图）"},
	{ID: "EnterPutrifiedDucts", Description: "进入腐殖渠（切图）"},
	{ID: "ToolPouch1", Description: "工具袋1（升级）"},
	{ID: "ToolPouch2", Description: "工具袋2（升级）"},
	{ID: "ToolPouch3", Description: "工具袋3（升级）"},
	{ID: "ToolPouch4", Description: "工具袋4（升级）"},
	{ID: "CraftingKit1", Description: "制作匣1（升级）"},
	{ID: "CraftingKit2", Description: "制作匣2（升级）"},
	{ID: "CraftingKit3", Description: "制作匣3（升级）"},
	{ID: "CraftingKit4", Description: "制作匣4（升级）"},
	{ID: "CurseCrest", Description: "被寄生（纹章）"},
	{ID: "GainedCurse", Description: "被寄生（事件）"},
	{ID: "WitchCrest", Description: "巫妪纹章（纹章）"},
	{ID: "WitchCrestTrans", Description: "巫妪纹章（切图）"},
	{ID: "Sylphsong", Description: "风灵谣（技能）"},
	{ID: "SylphsongTrans", Description: "风灵谣（切图）"},
	{ID: "ElegyOfTheDeep", Description: "深邃挽歌（技能）"},
	{ID: "GurrTheOutcastEncountered", Description: "被放逐的格尔-遭遇（Boss）"},
	{ID: "GurrTheOutcast", Description: "被放逐的格尔（Boss）"},
	{ID: "EnterWispThicket", Description: "进入火灵竹丛（切图）"},
	{ID: "EnterFatherOfTheFlame", Description: "进入炽焰之父（切图）"},
	{ID: "FatherOfTheFlame", Description: "炽焰之父（Boss）"},
	{ID: "LastJudgeEncountered", Description: "末日裁决者-遭遇（Boss）"},
	{ID: "MistCrossing", Description: "迷雾-中间长廊（切图）"},
	{ID: "EnterTheSlab", Description: "进入罪石牢狱-正门（切图）"},
	{ID: "SlabKeyIndolent", Description: "懒惰之钥（道具）"},
	{ID: "SlabKeyHeretic", Description: "异端之钥（道具）"},
	{ID: "SlabKeyApostate", Description: "叛教之钥（道具）"},
	{ID: "EnterFirstSinner", Description: "进入原罪者（切图）"},
	{ID: "FirstSinnerEncountered", Description: "原罪者-遭遇（Boss）"},
	{ID: "FirstSinner", Description: "原罪者（Boss）"},
	{ID: "EnterMountFay", Description: "进入费耶山（切图）"},
	{ID: "EnterBrightvein", Description: "进入冰晶脉窟（切图）"},
	{ID: "UpperMountFayTrans", Description: "费耶山-上层（切图）"},
	{ID: "EnterCogworkDancers", Description: "进入机枢舞者（切图）"},
	{ID: "CogworkDancersEncountered", Description: "机枢舞者-遭遇（Boss）"},
	{ID: "EnterRedMemory", Description: "进入赤红忆境（切图）"},
	{ID: "EnterDestroyedCogworks", Description: "进入被破坏的机枢核心（切图）"},
	{ID: "DestroyedCogworksVoidArena", Description: "被破坏的机枢核心遭遇战（小Boss）"},
	{ID: "DivingBellAbyssTrans", Description: "进入深渊-潜钟（切图）"},
	{ID: "EnterAbyss", Description: "进入深渊（切图）"},
	{ID: "LastDiveTrans", Description: "最终下潜（切图）"},
	{ID: "LostLaceEncountered", Description: "失心蕾丝-遭遇（Boss）"},
	{ID: "EnterWeavenestAtla", Description: "进入阿特拉织巢（切图）"},
	{ID: "EnterHalfwayHomeBasement", Description: "进入中途酒馆地下室（切图）"},
	{ID: "Lugoli", Description: "失格大厨（Boss）"},
	{ID: "PhantomTrans", Description: "幽影（切图）"},
	{ID: "TrailsEndTrans", Description: "沙克拉导师（切图）"},
	{ID: "Broodmother", Description: "育母（Boss）"},
	{ID: "EnterSandsOfKarak", Description: "进入卡拉卡沙川（切图）"},
	{ID: "EnterVoltnest", Description: "进入电荷巢穴（切图）"},
	{ID: "RagingConchfly", Description: "狂暴的螺蝇（Boss）"},
	{ID: "RagingConchflyTrans", Description: "狂暴的螺蝇（切图）"},
	{ID: "WatcherAtTheEdge", Description: "边陲守望者（Boss）"},
	{ID: "EnterSongclave", Description: "进入圣歌盟地（切图）"},
	{ID: "MetMergwin", Description: "梅尔格温（NPC）"},
	{ID: "GivenCouriersRasher", Description: "送货员熏肉交付（事件）"},
	{ID: "EnterWhiteward", Description: "进入白愈厅（切图）"},
	{ID: "PostWhitewardElevatorTrans", Description: "穿过白愈厅电梯后（切图）"},
	{ID: "CollectedSurgeonsKey", Description: "医师的钥匙（道具）"},
	{ID: "TheUnravelledEncountered", Description: "散茧魂渊-遭遇（Boss）"},
	{ID: "TheUnravelled", Description: "散茧魂渊（Boss）"},
	{ID: "SecondSentinelAwoken", Description: "唤醒次席戍卫（事件）"},
	{ID: "EnterSecondSentinel", Description: "进入次席戍卫（切图）"},
	{ID: "SecondSentinelBossEncountered", Description: "次席戍卫-遭遇（Boss）"},
	{ID: "SecondSentinel", Description: "次席戍卫（Boss）"},
	{ID: "FleaFestivalBegin", Description: "跳蚤庆典开始（事件）"},
	{ID: "FleaFestivalEnd", Description: "跳蚤庆典结束（事件）"},
	{ID: "NylethEncountered", Description: "尼莱斯-遭遇（Boss）"},
	{ID: "Nyleth", Description: "尼莱斯（Boss)"},
	{ID: "KhannEncountered", Description: "壳王卡汗-遭遇（Boss）"},
	{ID: "CrustKingKhann", Description: "壳王卡汗（Boss）"},
	{ID: "SkarrsingerKarmelita", Description: "斯卡尔歌后卡梅莉塔（Boss）"},
	{ID: "Palestag", Description: "苍白苜鹿（Boss）"},
	{ID: "CloverDancersEncountered", Description: "三叶草舞者-遭遇（Boss）"},
	{ID: "CloverDancers", Description: "三叶草舞者（Boss）"},
}

var cacheAliases = map[string]string{
	"Conchflies1":                   "GreatConchflies",
	"WhisperingVaultsGauntlet":      "WhisperingVaultsArena",
	"EnterHighHallsGauntlet":        "EnterHighHallsArena",
	"HighHallsGauntlet":             "HighHallsArena",
	"CollectedWhiteWardKey":         "CollectedWhitewardKey",
	"SavedFleaUnderworksExplosions": "SavedFleaUnderworksCauldron",
	"MountFayStation":               "SlabStation",
	"CityBellwayTube":               "GrandBellwayTube",
	"CollectedHeartNyleth":          "HeartNyleth",
	"CollectedHeartKhann":           "HeartKhann",
	"CollectedHeartKarmelita":       "HeartKarmelita",
	"CollectedHeartClover":          "HeartClover",
	"CompletedRedMemory":            "RedMemory",
}
