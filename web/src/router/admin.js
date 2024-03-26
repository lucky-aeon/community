// admin page

import Article from "@/views/admin/article/Article.vue";
import Code from "@/views/admin/code/Code.vue";
import File from "@/views/admin/file/File.vue";
import LoginLog from "@/views/admin/log/LoginLog.vue";
import OperLog from "@/views/admin/log/OperLog.vue";
import Member from "@/views/admin/member/Member.vue";
import Type from "@/views/admin/type/Type.vue";
import Tag from "@/views/admin/user/Tag.vue";
import UserList from "@/views/admin/user/UserList.vue";
import MessageTemplate from "@/views/admin/message/MessageTemplate.vue";

const AdminRouters = [
    {
        path: "member",
        component: Member
    },
    {
        path: "file",
        component: File
    },
    {
        path: "code",
        component: Code
    },
    {
        path: "oper/log",
        component: OperLog
    },
    {
        path: "login/log",
        component: LoginLog
    },
    {
        path: "type",
        component: Type
    },
    {
        path: "user",
        component: UserList
    },
    {
        path: "article",
        component: Article
    },
    {
        path: "user/tag",
        component: Tag
    },
    {
        path: "message/template",
        component: MessageTemplate
    }

]

export default AdminRouters