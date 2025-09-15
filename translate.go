package main

type SplitData struct {
	ID          string
	Description string
}

func getSplitDescriptionByID(id string) string {
	for _, split := range splitsCache {
		if split.ID == id {
			return split.Description
		}
	}
	return ""
}

func getIndexByID(id string) int {
	for i, split := range splitsCache {
		if split.ID == id {
			return i
		}
	}
	return -1
}

func getIDByDescription(desc string) string {
	for _, split := range splitsCache {
		if split.Description == desc {
			return split.ID
		}
	}
	return ""
}

var splitsCache = []*SplitData{
	{"ManualSplit", "手动分割（其它）"},
	{"StartNewGame", "开始新游戏（开始）"},
	{"EndingSplit", "任意结束（结束）"},
	{"EndingA", "苍白之母结局（结束）"},
	{"Menu", "主菜单（菜单）"},
	{"PlayerDeath", "死亡（事件）"},
	{"AnyTransition", "任意切图（切图）"},
	{"MossMother", "苔藓之母（Boss）"},
	{"SilkSpear", "丝之矛（技能）"},
	{"BellBeast", "钟道兽（Boss）"},
	{"SwiftStep", "冲刺（技能）"},
	{"Lace1", "蕾丝（Boss）"},
	{"DriftersCloak", "斗篷（技能）"},
	{"FourthChorus", "第四圣咏团（Boss）"},
	{"EnterGreymoor", "进入灰沼（切图）"},
	{"Moorwing", "荒沼翼主（Boss）"},
	{"ClingGrip", "爬墙（技能）"},
	{"Widow", "黑寡妇（Boss）"},
	{"EnterMist", "进入迷雾（切图）"},
	{"LeaveMist", "离开迷雾（切图）"},
	{"Phantom", "幽影（Boss）"},
	{"Act2Started", "第二幕开始（事件）"},
	{"CogworkDancers", "机枢舞者（Boss）"},
	{"Trobbio", "特罗比奥（Boss）"},
	{"Clawline", "飞针冲刺（技能）"},
	{"EnterHighHalls", "进入高庭（切图）"},
	{"EnterHighHallsGauntlet", "进入高庭遭遇战（切图）"},
	{"HighHallsGauntlet", "高庭遭遇战（小Boss）"},
	{"Lace2", "蕾丝2（Boss）"},
	{"VaultkeepersMelody", "管理员旋律（旋律）"},
	{"ArchitectsMelody", "建筑师旋律（旋律）"},
	{"ConductorsMelody", "指挥家旋律（旋律）"},
	{"UnlockedMelodyLift", "解锁三重旋律电梯（事件）"},
	{"NeedleUpgrade1", "织针升级1（升级）"},
	{"NeedleUpgrade2", "织针升级2（升级）"},
	{"NeedleUpgrade3", "织针升级3（升级）"},
	{"NeedleUpgrade4", "织针升级4（升级）"},
}
