import axios from 'axios';
import { ElMessage, ElMessageBox } from 'element-plus';
import { useUserStore } from '@/pinia/modules/user';
import router from '@/router/index';
import { showLoading, closeLoading } from '@/utils/loadingManager'; // 引入加载动画管理器
import requestManager from '@/utils/requestManager'; // 引入请求管理器

const service = axios.create({
    baseURL: import.meta.env.VITE_BASE_API,
    timeout: 99999,
});

// http request 拦截器
service.interceptors.request.use(
    config => {
        if (!config.donNotShowLoading) {
            showLoading(config.loadingOption);
        }

        // 使用 AbortController 管理请求
        config.signal = requestManager.addRequest();

        const userStore = useUserStore();
        config.headers = {
            'Content-Type': 'application/json',
            'x-token': userStore.token,
            'x-user-id': userStore.userInfo.ID,
            ...config.headers,
        };
        return config;
    },
    error => {
        if (!error.config.donNotShowLoading) {
            closeLoading();
        }
        ElMessage({
            showClose: true,
            message: error,
            type: 'error',
        });
        return Promise.reject(error);
    }
);

// http response 拦截器
service.interceptors.response.use(
    response => {
        if (!response.config.donNotShowLoading) {
            closeLoading();
        }

        const userStore = useUserStore();
        if (response.headers['new-token']) {
            userStore.setToken(response.headers['new-token']);
        }
        if (response.data.code === 0 || response.headers.success === 'true') {
            if (response.headers.msg) {
                response.data.msg = decodeURI(response.headers.msg);
            }
            return response.data;
        } else {
            ElMessage({
                showClose: true,
                message: response.data.msg || decodeURI(response.headers.msg),
                type: 'error',
            });
            return response.data.msg ? response.data : response;
        }
    },
    error => {
        if (!error.config.donNotShowLoading) {
            closeLoading();
        }

        if (error.name === 'AbortError') {
            console.log('请求已被取消');
            return Promise.reject(error); // 忽略取消的请求
        }

        if (!error.response) {
            // ElMessageBox.confirm(
            //     `<p>检测到请求错误</p><p>${error}</p>`,
            //     '请求报错',
            //     {
            //         dangerouslyUseHTMLString: true,
            //         distinguishCancelAndClose: true,
            //         confirmButtonText: '稍后重试',
            //         cancelButtonText: '取消',
            //     }
            // );
            return Promise.reject(error);
        }

        switch (error.response.status) {
            case 500:
            case 404:
            case 401:
                // 处理特定状态码...
                break;
        }

        return Promise.reject(error);
    }
);

export default service;