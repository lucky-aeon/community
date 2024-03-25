<template>
    <a-menu style="font-size: 16px;" class="line" mode="horizontal" :default-selected-keys="['1']">
        <a-menu-item key="0" :style="{ padding: 0, marginRight: '38px' }" disabled>

        </a-menu-item>
        <a-menu-item key="1">首页</a-menu-item>
      <a-menu-item key="2"><a href="https://xhyqaq.github.io/">付费</a></a-menu-item>
    </a-menu>
    <a-divider />
  <a-modal v-model:visible="visible" @ok="handleOk" @cancel="handleCancel" width="auto">
    <template #title>
      用户协议
    </template>
    <div>
      <h1>用户协议</h1>

      <h2>1. 账户注册</h2>
      <p>
        1.1 用户应提供真实、准确、完整的个人信息进行账户注册，并保证随时更新这些信息以确保其准确性。<br>
        1.2 用户注册成功后，将获得一个账户及相关登录凭证，用户应对账户及登录凭证的安全、保密负有责任。<br>
        1.3 您不得将账号借给他人使用，如有违反将导致账号使用权被收回。
      </p>

      <h2>2. 用户行为规范</h2>
      <p>
        2.1 用户在平台上的行为应符合国家法律法规、社会公德和平台规定的要求，用户不得利用平台服务从事以下行为：<br>
        - 发布、传播、储存违法、有害、侵权、淫秽等信息内容；<br>
        - 损害他人合法权益的行为，包括但不限于诽谤、侵犯隐私、恶意诋毁等；<br>
        - 传播垃圾信息、广告信息或其他未经授权的商业行为；<br>
        - 从事任何可能干扰、破坏平台正常运行的行为，包括但不限于病毒攻击、网络攻击等；<br>
        - 以任何方式干扰其他用户或平台管理人员的正常使用权益。<br>
        2.2 用户不得利用平台从事任何违法犯罪行为或违反公序良俗的行为，如散布谣言、制作、复制、传播淫秽物品等。
      </p>

      <h2>3. 付费服务</h2>
      <p>
        3.1 平台可能提供一些付费服务，用户可以选择购买并享受相应的服务内容。<br>
        3.2 当前收费为永久，然而请注意，平台保留根据服务提供的不同，将收费模式改为年费制的权利。若发生收费模式变更，平台将提前通知用户，并按照通知的规定进行收费。<br>
        3.3 用户购买付费服务的行为需符合平台发布的购买流程和规定，用户应按照平台要求支付相应的费用。<br>
        3.4 用户购买的付费服务在购买成功后将立即生效，用户无权要求退款，除非平台另有明确规定。
      </p>
      <h2>4. 免责声明</h2>
      <p>
        4.1 平台仅提供相关服务的技术支持，不对用户使用服务的结果承担任何责任。<br>
        4.2 用户在使用平台服务过程中应自行承担风险，平台不对因使用服务而导致的任何损失或损害承担责任。<br>
        4.3 用户应理解并接受，平台无法对其他用户的行为进行控制，因此平台对其他用户的行为及其可能导致的损失或损害不承担责任。<br>
        4.4 用户应妥善保管账户及登录凭证，对于因用户的疏忽、遗失或泄露导致的损失或损害，平台不承担责任。<br>
        4.5 平台有权根据实际情况对用户协议进行修改和更新，用户在使用平台服务时应及时查阅最新版本的用户协议。
      </p>
      <h2>5.知识产权</h2>
      <p>
        5.1 平台上的所有内容，包括但不限于文字、图形、标识、图像、音频、视频等知识产权，均归平台或相关权利人所有。<br>
        5.2 未经平台或相关权利人的明确许可，用户不得以任何形式进行复制、修改、展示、传播、发行或利用上述知识产权。<br>
        5.3 用户在平台上发布的内容，应确保拥有相关知识产权或获得了合法授权，否则将承担相关法律责任。
      </p>
      <h2>6.隐私保护</h2>
      <p>
        6.1 平台将按照相关法律法规和隐私政策对用户的个人信息进行保护。<br>
        6.2 平台将采取合理的安全措施保护用户的个人信息，防止未经授权的访问、使用或泄露。<br>
        6.3 用户应妥善保管个人信息，不得将个人信息泄露给其他第三方，如因用户自身原因导致个人信息泄露，平台不承担责任。
      </p>
      <h2>7.终止协议</h2>
      <p>
        7.1 平台有权在发现用户违反用户协议的情况下，立即终止与该用户的协议关系，并采取相应的法律措施。<br>
        7.2 协议终止后，用户应停止使用平台服务，并履行其在协议有效期间的相关义务。
      </p>
      <h2>8.法律适用和争议解决</h2>
      <p>
        8.1 用户协议的解释和适用以中华人民共和国法律为准据。<br>
        8.2 用户与平台之间因用户协议产生的争议，应通过友好协商解决；协商不成的，应提交有管辖权的人民法院解决。
      </p>
      <p>以上为用户协议的基本内容，如有任何疑问，请联系我们。</p>
    </div>
  </a-modal>
    <div style="height: 100%; width: 100%;margin-top: 100px;">
        <a-card hoverable style="margin: auto; width: 550px;" title="身份验证">
            <template #extra>
                <a-link @click="handleClick">用户协议</a-link>
            </template>
            <a-form :model="authForm" :style="{ width: '100%' }" @submit="handleSubmit">
                <a-form-item feedback field="account" tooltip="联系站长获取邀请码" label="邮箱"
                    :rules="[{ required: true, message: 'account is required' }, { minLength: 6, message: 'must be greater than 6 characters' }]"
                    :validate-trigger="['change', 'input']">
                    <a-input v-model="authForm.account" placeholder="please enter your emial..." />
                </a-form-item>
                <a-form-item feedback field="password" label="密码"
                    :rules="[{ required: true, message: 'password is required' }, { minLength: 6, message: 'must be greater than 6 characters' }]"
                    :validate-trigger="['change', 'input']">
                    <a-input v-model="authForm.password" type="password" placeholder="please enter your password..." />
                </a-form-item>
                <a-form-item field="code" tooltip="第一次登录需填写" label="邀请码">
                    <a-verification-code :formatter="(inputValue) => /^\d*$/.test(inputValue) ? inputValue : false"
                        :length="8" v-model="authForm.code" style="width: 200px" />
                </a-form-item>
                <a-form-item feedback v-if="authForm.code.length == 8" field="name" label="昵称"
                    :rules="[{ required: true, message: 'name is required' }, { minLength: 2, message: 'must be greater than 2 characters' }]"
                    :validate-trigger="['change', 'input']">
                    <a-input v-model="authForm.name" placeholder="please enter your name..." />
                </a-form-item>
                <a-form-item feedback field="isRead" :rules="[{ type:'boolean', true:true, message: '请阅读用户协议并同意' }]"
                    :validate-trigger="['change', 'input']">
                    <a-checkbox v-model="authForm.isRead"> I have read the manual </a-checkbox>
                </a-form-item>

                <a-form-item>
                    <a-button type="primary" html-type="submit">login / register</a-button>
                </a-form-item>
            </a-form>
        </a-card>
    </div>
</template>
<script setup>
import { useUserStore } from '@/stores/UserStore';
import {reactive, ref} from 'vue';

const authForm = reactive({
    account: "",
    password: "",
    code: "",
    name: "",
    isRead: false
})
const userStore = useUserStore()
const handleSubmit = () => {
    userStore.login(authForm)
}

const visible = ref(false);

const handleClick = () => {
  visible.value = true;
};
const handleOk = () => {
  visible.value = false;
};
const handleCancel = () => {
  visible.value = false;
}
</script>