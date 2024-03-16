// admin page

import Code from "@/views/admin/code/Code.vue";
import File from "@/views/admin/file/File.vue";
import Member from "@/views/admin/member/Member.vue";

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
    }
]

export default AdminRouters