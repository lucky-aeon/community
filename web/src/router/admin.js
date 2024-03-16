// admin page

import Code from "@/views/admin/code/Code.vue";
import File from "@/views/admin/file/File.vue";
import Member from "@/views/admin/member/Member.vue";
import Log from "@/views/admin/log/Log.vue";

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
        path: "log",
        component: Log
    }
]

export default AdminRouters