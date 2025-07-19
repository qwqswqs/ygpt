// utils/requestManager.js
const requestManager = {
    // 存储所有请求的 AbortController
    controllers: [],

    // 添加请求
    addRequest() {
        const controller = new AbortController();
        this.controllers.push(controller);
        return controller.signal;
    },

    // 取消所有请求
    cancelAll() {
        this.controllers.forEach(controller => {
            controller.abort('页面切换，请求被取消');
        });
        this.controllers = []; // 清空控制器列表
    },
};

export default requestManager;