<template>
    <el-alert title="如果想为本项目做贡献，请前往本项目的仓库地址： https://github.com/CuteReimu/sssplitmaker" type="info" effect="dark" close-text="前往" @close="openGithub" style="max-width: 960px"></el-alert>
    <el-select v-model="currentTemplate" filterable placeholder="你可以选择现有模板" style="width: 500px" @change="selectTemplate">
        <el-option v-for="item in templates" :key="item.value" :label="item.label" :value="item.value"></el-option>
    </el-select>
    <el-upload drag accept=".lss" :auto-upload="false" :show-file-list="false" :on-change="handleChange">
        <el-icon class="el-icon--upload"><upload-filled style="width: 80px;"></upload-filled></el-icon>
        <div class="el-upload__text">
            你也可以将文件拖拽到这里或者 <em>点击上传</em> 只支持 *.lss 文件
        </div>
    </el-upload>
    <div style="display: flex; gap: 8px;">
        <el-button type="success" @click="fillIcons">一键填充所有未填充的图标</el-button>
        <el-button type="danger" @click="resetIcons">一键清空所有图标</el-button>
        <el-button type="primary" @click="downloadIcons">下载所有图标</el-button>
        <el-text style="margin: 0px 10px;">Auto Splitter版本：1.25.4</el-text>
        <el-button type="warning" @click="fixLiveSplit" :loading="fixingLiveSplit">更新LiveSplit</el-button>
        <el-button type="warning" @click="openHelp">查看帮助</el-button>
    </div>
    <div>
        <el-switch
            v-model="skipStartAnimation"
            size="large"
            active-text="跳过开场动画"
            inactive-text="不跳过"
            @change="onSkipStartAnimationChange"
            :disabled="disableStartAnimation"
        ></el-switch>
        <el-text size="large" style="margin-left: 12px; cursor: pointer; text-decoration: underline;" @click="openSkipAnimationHelp">这是什么意思？</el-text>
    </div>
    <el-table :data="tableData" style="max-width: 960px">
        <el-table-column label="图标" width="60px">
            <template #default="scope">
                <el-image v-if="scope.row.icon.length>0" style="width: 25px; height: 25px" :src="scope.row.icon" fit="contain"></el-image>
            </template>
        </el-table-column>
        <el-table-column label="节点名称">
            <template #default="scope">
                <el-input v-if="scope.$index>0 && scope.$index<tableData.length-1" v-model="scope.row.name" placeholder="节点名称" style="width: 300px"></el-input>
            </template>
        </el-table-column>
        <el-table-column label="触发事件">
            <template #default="scope">
                <el-select-v2 v-if="scope.$index<tableData.length-1" v-model="scope.row.event" :options="options"
                           @change="onEventChange(scope.$index)" filterable placeholder="触发事件" style="width: 300px" />
            </template>
        </el-table-column>
        <el-table-column label="操作" :width="220">
            <template #default="scope">
                <el-button v-if="scope.$index>0" :icon="Plus" circle @click="addLine(scope.$index)"></el-button>
                <el-button :disabled="tableData.length<=3" v-if="scope.$index>0 && scope.$index<tableData.length-1" :icon="Minus" circle @click="removeLine(scope.$index)"></el-button>
                <el-button :disabled="scope.$index<=1" v-if="scope.$index>0 && scope.$index<tableData.length-1" :icon="Top" circle @click="swapLine(scope.$index-1, scope.$index)"></el-button>
                <el-button :disabled="scope.$index>=tableData.length-2" v-if="scope.$index>0 && scope.$index<tableData.length-1" :icon="Bottom" @click="swapLine(scope.$index, scope.$index+1)" circle></el-button>
            </template>
        </el-table-column>
    </el-table>
    <div>
        <el-button type="primary" @click="submit" style="align-self: flex-start;" :disabled="disableSubmit">另存为</el-button>
        <el-checkbox v-model="includeTimeRecords" size="large" style="margin-left: 20px">保留*.lss文件中原本的时间记录（如果有）</el-checkbox>
    </div>
</template>

<script setup lang="ts">
import { ElAlert, ElSelect, ElSelectV2, ElOption, ElUpload, ElButton, ElSwitch, ElTable, ElTableColumn, ElCheckbox, ElMessage, ElText, ElIcon, ElImage, ElInput } from 'element-plus';
import { Plus, Minus, Top, Bottom, UploadFilled } from '@element-plus/icons-vue';
import { ref, onMounted } from 'vue';
import { GetOptions, GetTemplates, LoadSplitFile, GetSplits, GetIcon, SaveSplitsFile, SaveIconsZip, FixLiveSplit } from '../wailsjs/go/main/App';
import { BrowserOpenURL, LogError, EventsOn } from '../wailsjs/runtime';

interface Row {
  name: string;
  event: string;
  icon: string;
  other?: unknown[];
}

interface Option {
  value: string;
  label: string;
}

const includeTimeRecords = ref(true);
const disableSubmit = ref(false);
const skipStartAnimation = ref(false);
const disableStartAnimation = ref(false);
const options = ref<Option[]>([]);
const currentTemplate = ref('');
const templates = ref<Option[]>([]);
const fixingLiveSplit = ref(false);
const tableData = ref<Row[]>([
  { name: '', event: 'StartNewGame', icon: '' },
  { name: '任意结束', event: 'EndingSplit', icon: '' },
  { name: '', event: 'ManualSplit', icon: '' },
]);

onMounted(async () => {
  try {
    options.value = await GetOptions();
  } catch (e) {
    LogError(e);
  }
  try {
    templates.value = await GetTemplates();
  } catch (e) {
    LogError(e);
  }
  try {
    tableData.value[1].icon = await GetIcon(tableData.value[1].event);
  } catch (e) {
    LogError(e);
  }
});

function refreshStartAnimationChange(eventValue: string) {
  switch (eventValue) {
    case 'StartNewGame':
      skipStartAnimation.value = false;
      disableStartAnimation.value = false;
      break;
    case 'Act1Start':
      skipStartAnimation.value = true;
      disableStartAnimation.value = false;
      break;
    default:
      disableStartAnimation.value = true;
  }
}

function onSkipStartAnimationChange(value: boolean) {
  if (value) {
    if (tableData.value[0].event !== 'Act1Start') tableData.value[0].event = 'Act1Start';
  } else {
    if (tableData.value[0].event !== 'StartNewGame') tableData.value[0].event = 'StartNewGame';
  }
}

function addLine(index: number) {
  tableData.value.splice(index, 0, { name: '手动分割', event: 'ManualSplit', icon: '' });
}

function removeLine(index: number) {
  tableData.value.splice(index, 1);
}

function swapLine(index1: number, index2: number) {
  const temp = tableData.value[index1];
  tableData.value[index1] = tableData.value[index2];
  tableData.value[index2] = temp;
}

async function submit() {
  disableSubmit.value = true;
  try {
    await SaveSplitsFile(tableData.value.slice(0, -1) as any, includeTimeRecords.value);
  } catch (e) {
    LogError(e);
    ElMessage({ message: '导出失败', type: 'error', plain: true });
  } finally {
    disableSubmit.value = false;
  }
}

async function downloadIcons() {
  disableSubmit.value = true;
  try {
    await SaveIconsZip();
  } catch (e) {
    LogError(e);
    ElMessage({ message: '导出失败', type: 'error', plain: true });
  } finally {
    disableSubmit.value = false;
  }
}

async function handleChange(file: { raw: File }) {
  if (!file?.raw) return;
  try {
    const text = await file.raw.text();
    const newData = await LoadSplitFile(text);
    tableData.value = newData as Row[];
    refreshStartAnimationChange(tableData.value[0]?.event ?? '');
  } catch (e) {
    LogError(e);
    ElMessage({ message: String(e), type: 'error', plain: true });
  }
}

async function selectTemplate(value: string) {
  try {
    const res = await GetSplits(value);
    tableData.value = [...res.splits as Row[], { name: '', event: 'ManualSplit', icon: '' }];
    refreshStartAnimationChange(tableData.value[0].event);
  } catch (e) {
    LogError(e);
  }
}

function openGithub() {
  BrowserOpenURL('https://github.com/CuteReimu/sssplitmaker');
}

function openSkipAnimationHelp() {
  BrowserOpenURL('https://cutereimu.cn/daily/silksong/speedrun-submit.html#%E6%96%B0%E8%A7%84%E5%88%99-%E5%85%81%E8%AE%B8%E8%B7%B3%E8%BF%87%E5%BC%80%E5%9C%BA%E5%8A%A8%E7%94%BB');
}

function openHelp() {
  BrowserOpenURL('https://cutereimu.cn/daily/silksong/sssplitmaker-faq.html');
}

async function onEventChange(idx: number) {
  const eventValue = tableData.value[idx].event;
  if (idx === 0) refreshStartAnimationChange(eventValue);
  const opt = options.value.find(o => o.value === eventValue);
  if (opt) {
    const pos = opt.label.indexOf('（');
    tableData.value[idx].name = pos === -1 ? opt.label : opt.label.slice(0, pos);
  }
  if (idx === 0) return;
  try {
    tableData.value[idx].icon = await GetIcon(tableData.value[idx].event);
  } catch (e) {
    LogError(e);
  }
}

async function fillIcons() {
  for (let idx = 1; idx < tableData.value.length - 1; idx++) {
    const row = tableData.value[idx];
    if (row.icon.length === 0) {
      try {
        row.icon = await GetIcon(row.event);
      } catch (e) {
        LogError(e);
      }
    }
  }
}

function resetIcons() {
  for (let idx = 1; idx < tableData.value.length - 1; idx++) {
    tableData.value[idx].icon = '';
  }
}

async function fixLiveSplit() {
  fixingLiveSplit.value = true;
  try {
    await FixLiveSplit();
  } finally {
    fixingLiveSplit.value = false;
  }
}

EventsOn("ElMessage", (type, message) => {
  ElMessage({ message, type, plain: true });
});
</script>
