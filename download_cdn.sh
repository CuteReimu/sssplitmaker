#!/bin/sh
mkdir -p static
curl https://unpkg.com/vue@3.5.22/dist/vue.global.prod.js -o static/vue.global.prod.js
curl https://unpkg.com/element-plus@2.11.3/dist/index.full.min.js -o static/index.full.min.js
curl https://unpkg.com/element-plus@2.11.3/dist/index.css -o static/index.css
curl https://unpkg.com/@element-plus/icons-vue@2.3.2/dist/index.iife.min.js -o static/index.iife.min.js