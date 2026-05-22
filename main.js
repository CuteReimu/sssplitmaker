// 导入需要的图标
const { Plus, Minus, Top, Bottom, UploadFilled } = ElementPlusIconsVue;
const { ElMessage } = ElementPlus;
const { createApp } = Vue;
const app = createApp({
    data() {
        return {
            includeTimeRecords: true,
            disableSubmit: false,
            skipStartAnimation: false,
            disableStartAnimation: false,
            options: [],
            currentTemplate: '',
                templates: [],
                loading: false,
                tableData: [
                {
                    name: '',
                    event: 'StartNewGame',
                    icon: '',
                },
                {
                    name: '任意结束',
                    event: 'EndingSplit',
                    icon: '',
                },
                {
                    name: '',
                    event: 'ManualSplit',
                    icon: '',
                },
            ],
            fileList: [],
        };
    },
    setup() {
        return {
            Plus,
            Minus,
            Top,
            Bottom,
            UploadFilled,
        };
    },
    mounted() {
        axios.get('/get-options')
            .then(response => {
                this.options = response.data;
            })
            .catch(error => {
                console.log(error);
            })
        axios.get('/get-templates')
            .then(response => {
                this.templates = response.data;
            })
            .catch(error => {
                console.log(error);
            })
    },
    methods: {
        onSkipStartAnimationChange(value) {
            if (value) {
                if (this.tableData[0].event !== "Act1Start") {
                    this.tableData[0].event = "Act1Start";
                }
            } else {
                if (this.tableData[0].event !== "StartNewGame") {
                    this.tableData[0].event = "StartNewGame";
                }
            }
        },
        refreshStartAnimationChange(eventValue) {
            switch (eventValue) {
                case "StartNewGame":
                    this.skipStartAnimation = false;
                    this.disableStartAnimation = false;
                    break;
                case "Act1Start":
                    this.skipStartAnimation = true;
                    this.disableStartAnimation = false;
                    break;
                default:
                    this.disableStartAnimation = true;
            }
        },
        addLine(index) {
            this.tableData.splice(index, 0,
                {
                    name: '手动分割',
                    event: 'ManualSplit',
                    icon: '',
                },
            );
        },
        removeLine(index) {
            this.tableData.splice(index, 1);
        },
        swapLine(index1, index2) {
            const temp = this.tableData[index1];
            this.tableData[index1] = this.tableData[index2];
            this.tableData[index2] = temp;
        },
        submit() {
            axios.post('/build-splits',
                { data: this.tableData.slice(0, -1), includeTimeRecords: this.includeTimeRecords },
                {
                    responseType: 'blob',
                    headers: { 'Content-Type': 'application/json' },
                })
                .then(response => {
                    this.disableSubmit = true;
                    const blob = new Blob([response.data], { type: 'application/octet-stream' });
                    const url = URL.createObjectURL(blob);
                    const link = document.createElement('a');
                    link.href = url;
                    link.download = 'splits.lss';
                    link.style.display = 'none';
                    document.body.appendChild(link);
                    link.click();
                    document.body.removeChild(link);
                    setTimeout(() => {
                        URL.revokeObjectURL(url);
                        this.disableSubmit = false;
                    }, 2000);
                })
                .catch(error => {
                    console.error(error);
                    ElMessage({ message: '导出失败', type: 'error', plain: true });
                });
        },
        downloadIcons() {
            axios.get('/download/icons',
                {
                    responseType: 'blob',
                })
                .then(response => {
                    this.disableSubmit = true;
                    const blob = new Blob([response.data], { type: 'application/zip' });
                    const url = URL.createObjectURL(blob);
                    const link = document.createElement('a');
                    link.href = url;
                    link.download = 'icons.zip';
                    link.style.display = 'none';
                    document.body.appendChild(link);
                    link.click();
                    document.body.removeChild(link);
                    setTimeout(() => {
                        URL.revokeObjectURL(url);
                        this.disableSubmit = false;
                    }, 2000);
                })
                .catch(error => {
                    console.error(error);
                    ElMessage({ message: '导出失败', type: 'error', plain: true });
                });
        },
        uploadTableData(newData) {
            this.loading = true;
            this.tableData = newData;
            this.$nextTick(() => {
                this.loading = false;
            });
        },
        uploadFailed(err) {
            console.log(err);
            ElMessage({
                message: err,
                type: 'error',
                plain: true,
            });
        },
        selectTemplate(value) {
            this.loading = true;
            axios.post('/get-splits', {name: value}, {
                headers: {
                    'Content-Type': 'application/x-www-form-urlencoded',
                },
            })
                .then(response => {
                    this.tableData = [...response.data.splits, {
                        name: '',
                        event: 'ManualSplit',
                        icon: '',
                    }];
                    this.refreshStartAnimationChange(this.tableData[0].event)
                    this.$nextTick(() => {
                        this.loading = false;
                    });
                })
                .catch(error => {
                    console.log(error);
                    this.loading = false;
                });
        },
        handleChange(file, fileList) {
            const latest = fileList.slice(-1);
            this.fileList.splice(0, this.fileList.length, ...latest);
        },
        openTranslate() {
            window.open('/translate', '_blank', 'noopener,noreferrer');
        },
        openGithub() {
            window.open('https://github.com/CuteReimu/sssplitmaker', '_blank', 'noopener,noreferrer');
        },
        onEventChange(idx) {
            const eventValue = this.tableData[idx].event;
            if (idx === 0) {
                this.refreshStartAnimationChange(eventValue)
            }
            const opt = this.options.find(o => o.value === eventValue);
            if (opt) {
                const label = opt.label;
                const pos = label.indexOf('（');
                this.tableData[idx].name = pos === -1 ? label : label.slice(0, pos);
            }
            if (idx === 0) {
                return;
            }
            axios.get('/get-icon?name=' + encodeURIComponent(this.tableData[idx].event))
                .then(response => {
                    this.tableData[idx].icon = response.data;
                })
                .catch(error => {
                    console.log(error);
                });
        },
        fillIcons() {
            for (let idx = 1; idx < this.tableData.length - 1; idx++) {
                const row = this.tableData[idx];
                if (row.icon.length === 0) {
                    axios.get('/get-icon?name=' + encodeURIComponent(row.event))
                        .then(response => {
                            row.icon = response.data;
                        })
                        .catch(error => {
                            console.log(error);
                        });
                }
            }
        },
        resetIcons() {
            for (let idx = 1; idx < this.tableData.length - 1; idx++) {
                if (idx === 0 || idx === this.tableData.length - 1) {
                    continue;
                }
                this.tableData[idx].icon = '';
            }
        },
    },
});
app.component(UploadFilled.name, UploadFilled);
app.use(ElementPlus);
app.mount('#app');