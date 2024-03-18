// admin page

import Code from "@/views/admin/code/Code.vue";
import File from "@/views/admin/file/File.vue";
import Member from "@/views/admin/member/Member.vue";
import OperLog from "@/views/admin/log/OperLog.vue";
import LoginLog from "@/views/admin/log/LoginLog.vue";
import Type from "@/views/admin/type/Type.vue";

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
    }

]

export default AdminRouters