// utils/loadingManager.js
import { ElLoading } from 'element-plus';

let loadingInstance = null; // 定义全局的加载实例
let activeAxios = 0; // 记录当前活动的请求数量
let timer = null; // 延迟显示加载动画的定时器

// 显示加载动画
export const showLoading = (option = { target: null }) => {
    const loadDom = document.getElementById('gva-base-load-dom');
    activeAxios++;
    if (timer) {
        clearTimeout(timer);
    }
    timer = setTimeout(() => {
        if (activeAxios > 0) {
            if (!option.target) option.target = loadDom;
            loadingInstance = ElLoading.service(option); // 初始化加载实例
        }
    }, 400);
};

// 关闭加载动画
export const closeLoading = () => {
    activeAxios--;
    if (activeAxios <= 0) {
        clearTimeout(timer);
        if (loadingInstance) {
            loadingInstance.close(); // 关闭加载动画
            loadingInstance = null; // 清空加载实例
        }
    }
};