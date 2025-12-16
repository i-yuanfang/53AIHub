<template>
  <el-dialog v-model="visible" :title="title" width="500px" :close-on-click-modal="false">
    <el-form ref="formRef" :model="form" label-width="100px">
      <template v-if="type === 'wecom_login'">
        <el-form-item label="CorpID" prop="corp_id">
          <el-input v-model="form.corp_id" placeholder="请输入企业微信 CorpID" />
        </el-form-item>
        <el-form-item label="AgentID" prop="agent_id">
          <el-input v-model="form.agent_id" placeholder="请输入企业微信 AgentID" />
        </el-form-item>
        <el-form-item label="Secret" prop="secret">
          <el-input v-model="form.secret" type="password" show-password placeholder="请输入企业微信 Secret" />
        </el-form-item>
      </template>
      <template v-else-if="type === 'dingtalk_login'">
        <el-form-item label="AppKey" prop="app_key">
          <el-input v-model="form.app_key" placeholder="请输入钉钉 AppKey" />
        </el-form-item>
        <el-form-item label="AppSecret" prop="app_secret">
          <el-input v-model="form.app_secret" type="password" show-password placeholder="请输入钉钉 AppSecret" />
        </el-form-item>
      </template>
      <template v-else-if="type === 'feishu_login'">
        <el-form-item label="AppID" prop="app_id">
          <el-input v-model="form.app_id" placeholder="请输入飞书 AppID" />
        </el-form-item>
        <el-form-item label="AppSecret" prop="app_secret">
          <el-input v-model="form.app_secret" type="password" show-password placeholder="请输入飞书 AppSecret" />
        </el-form-item>
      </template>
      <el-form-item label="启用">
        <el-switch v-model="enabled" />
      </el-form-item>
    </el-form>
    <template #footer>
      <span class="dialog-footer">
        <el-button @click="visible = false">取消</el-button>
        <el-button type="primary" @click="handleSave" :loading="loading">保存</el-button>
      </span>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue'
import { ElMessage } from 'element-plus'
import enterpriseConfigApi from '@/api/modules/enterprise-config'

const props = defineProps<{
  onSuccess?: () => void
}>()

const visible = ref(false)
const loading = ref(false)
const type = ref('')
const enabled = ref(false)
const form = ref<any>({})

const title = computed(() => {
  switch (type.value) {
    case 'wecom_login':
      return '企业微信配置'
    case 'dingtalk_login':
      return '钉钉配置'
    case 'feishu_login':
      return '飞书配置'
    default:
      return '配置'
  }
})

const open = async (configType: string) => {
  type.value = configType
  visible.value = true
  loading.value = true
  try {
    const res = await enterpriseConfigApi.get(configType)
    if (res.data) {
      enabled.value = res.data.enabled
      try {
        form.value = JSON.parse(res.data.content)
      } catch (e) {
        form.value = {}
      }
    }
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

const handleSave = async () => {
  loading.value = true
  try {
    await enterpriseConfigApi.save(type.value, {
      content: JSON.stringify(form.value),
      enabled: enabled.value
    })
    ElMessage.success('保存成功')
    visible.value = false
    props.onSuccess?.()
  } catch (error) {
    console.error(error)
  } finally {
    loading.value = false
  }
}

defineExpose({
  open
})
</script>
