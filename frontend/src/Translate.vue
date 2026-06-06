<template>
    <el-table :data="tableData" border stripe v-loading="loading" element-loading-text="正在加载数据...">
        <el-table-column prop="description" label="Description" :width="350"></el-table-column>
        <el-table-column prop="translate" label="翻译" :width="250"></el-table-column>
        <el-table-column prop="tooltip" label="Tooltip"></el-table-column>
        <el-table-column label="Key" :width="310">
            <template #default="scope">
                <div>{{ scope.row.key }}</div>
                <div v-if="scope.row.alias">{{ scope.row.alias }}</div>
            </template>
        </el-table-column>
    </el-table>
</template>

<script setup lang="ts">
import { ElTable, ElTableColumn } from 'element-plus';
import { ref, onMounted } from 'vue';
import { GetTranslateData } from '../wailsjs/go/main/App';

interface TranslateRow {
  key: string;
  translate: string;
  description: string;
  tooltip: string;
  alias?: string | string[];
}

const loading = ref(true);
const tableData = ref<TranslateRow[]>([]);

onMounted(async () => {
  try {
    const json = await GetTranslateData();
    tableData.value = JSON.parse(json);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
});
</script>
