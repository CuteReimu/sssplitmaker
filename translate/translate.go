package translate

type SplitData struct {
	ID          string
	Description string
}

func GetSplitDescriptionByID(id string) string {
	for _, split := range SplitsCache {
		if split.ID == id {
			return split.Description
		}
	}
	return ""
}

func GetIndexByID(id string) int {
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

var SplitsCache = []*SplitData{
	{"ManualSplit", "手动分割（其它）"},
	{"StartNewGame", "开始新游戏（开始）"},
	{"EndingSplit", "任意结束（结束）"},
	{"EndingA", "苍白之母结局（结束）"},
	{"Menu", "主菜单（菜单）"},
	{"PlayerDeath", "死亡（事件）"},
	{"AnyTransition", "任意切图（切图）"},
	{"MossMother", "苔藓之母（Boss）"},
	{"MossMotherTrans", "苔藓之母（切图）"},
	{"SilkSpear", "丝之矛（技能）"},
	{"SilkSpearTrans", "丝之矛（切图）"},
	{"BellBeast", "钟道兽（Boss）"},
	{"BellBeastTrans", "钟道兽（切图）"},
	{"MarrowBell", "髓骨洞窟敲钟（事件）"},
	{"SwiftStep", "冲刺（技能）"},
	{"SwiftStepTrans", "冲刺（切图）"},
	{"Lace1", "蕾丝（Boss）"},
	{"Lace1Trans", "蕾丝（切图）"},
	{"DeepDocksBell", "深坞敲钟（事件）"},
	{"DriftersCloak", "斗篷（技能）"},
	{"DriftersCloakTrans", "斗篷（切图）"},
	{"FourthChorus", "第四圣咏团（Boss）"},
	{"EnterGreymoor", "进入灰沼（切图）"},
	{"GreymoorBell", "灰沼敲钟（事件）"},
	{"Moorwing", "荒沼翼主（Boss）"},
	{"MoorwingTrans", "荒沼翼主（切图）"},
	{"ClingGrip", "爬墙（技能）"},
	{"ClingGripTrans", "爬墙（切图）"},
	{"ShellwoodBell", "甲木林敲钟（事件）"},
	{"Widow", "黑寡妇（Boss）"},
	{"BellhartBell", "钟心镇敲钟（事件）"},
	{"LastJudge", "末代裁决者（Boss）"},
	{"EnterMist", "进入迷雾（切图）"},
	{"LeaveMist", "离开迷雾（切图）"},
	{"Phantom", "幽影（Boss）"},
	{"Act2Started", "第二幕开始（事件）"},
	{"CogworkDancers", "机枢舞者（Boss）"},
	{"WhisperingVaultsGauntlet", "低语书库遭遇战（小Boss）"},
	{"Trobbio", "特罗比奥（Boss）"},
	{"TrobbioTrans", "特罗比奥（切图）"},
	{"Clawline", "飞针冲刺（技能）"},
	{"EnterHighHalls", "进入高庭（切图）"},
	{"EnterHighHallsGauntlet", "进入高庭遭遇战（切图）"},
	{"HighHallsGauntlet", "高庭遭遇战（小Boss）"},
	{"Lace2", "蕾丝2（Boss）"},
	{"VaultkeepersMelody", "管理员旋律（旋律）"},
	{"VaultkeepersMelodyTrans", "管理员旋律（切图）"},
	{"ArchitectsMelody", "建筑师旋律（旋律）"},
	{"ArchitectsMelodyTrans", "建筑师旋律（切图）"},
	{"ConductorsMelody", "指挥家旋律（旋律）"},
	{"ConductorsMelodyTrans", "指挥家旋律（切图）"},
	{"UnlockedMelodyLift", "解锁三重旋律电梯（事件）"},
	{"NeedleUpgrade1", "织针升级1（升级）"},
	{"NeedleUpgrade2", "织针升级2（升级）"},
	{"NeedleUpgrade3", "织针升级3（升级）"},
	{"NeedleUpgrade4", "织针升级4（升级）"},
}
