import { apiUserTagArticleCount } from '@/apis/tags'
import { apiAuthAccount, getUserInfo, getUserMenu } from '@/apis/user'
import router from '@/router'
import * as authUtls from '@/utils/auth'
import { defineStore } from 'pinia'

// 你可以对 `defineStore()` 的返回值进行任意命名，但最好使用 store 的名字，同时以 `use` 开头且以 `Store` 结尾。(比如 `useUserStore`，`useCartStore`，`useProductStore`)
// 第一个参数是你的应用中 Store 的唯一 ID。
export const useUserStore = defineStore('UserStore', {
  // state
  state: () => ({
    token: "",
    userInfo: {},
    userTags: [],
    menu: []
  }),

  // actions
  actions: {
    isLogin() {
      if (!authUtls.isLogin()) {
        router.push("/auth")
        return false
      }
      return true
    },
    async refreshTags() {
      if (!this.isLogin()) return
      apiUserTagArticleCount().then(({ data, ok }) => {
        if (!ok) {
          this.userTags = []
          return
        }
        this.userTags = data
      })
    },
    async refreshInfo() {
      if (!this.isLogin()) return
      this.refreshTags()
      const { data, ok } = await getUserInfo()
      if (!ok) {
        this.logOut()
        return
      }
      this.userInfo = data
    },
    async logOut() {
      authUtls.clearToken()
      this.token = undefined
      this.userInfo = {}
      this.menu = []
      router.push("/auth")
    },
    async login(authForm) {
      authUtls.clearToken()
      if (!authForm.isRead) {
        return
      }
      let temp = Object.assign({}, authForm)
      if (authForm.code.length != 8) {
        temp.code = 0
      }
      const { data, ok } = await apiAuthAccount(temp)
      if (!ok) {
        return
      }
      // 登录成功，保存token
      this.token = data.token
      authUtls.setToken(this.token)
      this.getMenu()
      this.refreshInfo()
      router.push("/")
    },
    async getMenu() {
      const { data, ok } = await getUserMenu()
      if (ok) this.menu = data
      else this.menu = []
    }
  }

  // getters

})