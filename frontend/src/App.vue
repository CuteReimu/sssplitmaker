<template>
    <div style="display: flex; flex-direction: column; align-items: center; gap: 12px; padding: 16px;">
        <div style="display: flex; gap: 8px;">
            <el-button :type="route.path === '/' ? 'primary' : ''" @click="router.push('/')">分割点生成</el-button>
            <el-button :type="route.path === '/translate' ? 'primary' : ''" @click="router.push('/translate')">翻译对照</el-button>
        </div>
        <router-view style="display: flex; flex-direction: column; align-items: center; gap: 12px; width: 100%;"></router-view>
    </div>
</template>

<script setup lang="ts">
import { ElButton } from 'element-plus';
import { useRouter, useRoute } from 'vue-router';

const router = useRouter();
const route = useRoute();
</script>