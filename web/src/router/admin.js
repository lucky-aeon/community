// admin page

import Code from "@/views/admin/code/Code.vue";
import File from "@/views/admin/file/File.vue";
import Member from "@/views/admin/member/Member.vue";
import OperLog from "@/views/admin/log/OperLog.vue";
import LoginLog from "@/views/admin/log/LoginLog.vue";
import Type from "@/views/admin/type/Type.vue";
import User from "@/views/admin/user/User.vue";
import Tag from "@/views/admin/user/Tag.vue";
import Article from "@/views/admin/article/Article.vue";

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
        component: User
    },
    {
        path: "article",
        component: Article
    },
    {
        path: "user/tag",
        component: Tag
    }

]

export default AdminRouters