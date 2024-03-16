import { isLogin } from "@/utils/auth";
import { useUserStore } from "./UserStore";

export default function usePermission() {
  const userStore = useUserStore();
  return {
    accessRouter(route) { // 判断当前用户是否有该路由的权限
      return (
        !route.meta?.requiresAuth ||
        !route.meta?.roles ||
        route.meta?.roles?.includes('*') ||
        (isLogin() && route.meta?.roles?.includes(userStore.userInfo.role))
      );
    },
    // You can add any rules you want
  };
}