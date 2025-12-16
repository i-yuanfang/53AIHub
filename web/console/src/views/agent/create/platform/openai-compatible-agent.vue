<template>
  <el-form ref="formRef" :model="store.platformConfig" :rules="rules" label-position="top">
    <el-form-item :label="t('Base URL')" prop="base_url">
      <el-input
        v-model="store.platformConfig.base_url"
        :placeholder="t('例如: https://api.openai.com')"
      />
    </el-form-item>

    <el-form-item :label="t('API Key')" prop="api_key">
      <el-input
        v-model="store.platformConfig.api_key"
        :placeholder="t('请输入 API Key')"
      />
    </el-form-item>

    <el-form-item :label="t('模型名称')" prop="model">
      <el-input
        v-model="store.platformConfig.model"
        :placeholder="t('例如 gpt-4, claude-3-opus...')"
      />
    </el-form-item>

  </el-form>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useCreateAgentStore } from '@/views/agent/create/store'
import type { FormInstance, FormRules } from 'element-plus'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()
const formRef = ref<FormInstance>()
const store = useCreateAgentStore()

// 初始化 platformConfig 中的字段
if (!store.platformConfig.base_url) {
  store.platformConfig.base_url = ''
}
if (!store.platformConfig.api_key) {
  store.platformConfig.api_key = ''
}
// ++ 初始化新字段 ++
if (!store.platformConfig.model) {
  store.platformConfig.model = ''
}

// ++ 为新字段添加规则 ++
const rules = ref<FormRules>({
  base_url: [
    { required: true, message: t('请输入 Base URL'), trigger: 'blur' }
  ],
  api_key: [
    { required: true, message: t('请输入 API Key'), trigger: 'blur' }
  ],
  model: [ // ++ 新增 ++
    { required: true, message: t('请输入模型名称'), trigger: 'blur' }
  ]
})

// ... (validate 和 resetField 方法保持不变) ...
const validate = () => {
  return formRef.value?.validate()
}

const resetField = () => {
  formRef.value?.resetFields()
}

defineExpose({
  validate,
  resetField
})
</script>

    <ElForm ref="agentFormRef" :model="store.form_data" label-width="104px" label-position="top">
      <template v-if="showChannelConfig">
        <div class="text-base text-[#1D1E1F] font-medium mt-6 mb-4">
          {{ $t('basic_info') }}
        </div>
        <AgentInfo v-model="store.form_data" />
      </template>
      <template v-else>
        <template v-if="store.agent_type === AGENT_TYPES.DIFY_WORKFLOW">
          <FieldInput
            v-model:list="store.form_data.settings.input_fields"
            :title="$t('agent.input_variable')"
            allow-update
            allow-add
            :update-request="inputUpdateRequest"
            type="input"
            :agent-type="store.agent_type"
          />
          <FieldInput
            v-model:list="store.form_data.settings.output_fields"
            :title="$t('agent.output_variable')"
            allow-add
            type="output"
            :agent-type="store.agent_type"
          />
          <RelateApp />
        </template>
        <template v-else>
          <BaseConfig />
          <RelateApp />
          <ExpandConfig />
        </template>
        <UseScope />
      </template>
    </ElForm>
  </div>
</template>

<script setup lang="ts">
import { inject, reactive, ref, watch } from 'vue'
import AgentInfo from '../components/agent-info.vue'
import BaseConfig from '../components/base-config.vue'
import UseScope from '../components/use-scope.vue'
import FieldInput from '../components/field-input.vue'
import AgentType from '../components/agent-type.vue'
import RelateApp from '../components/relate-agents.vue'

import { useAgentFormStore } from '../store'
import { generateInputRules } from '@/utils/form-rule'
import { generateRandomId } from '@/utils'
import md5 from '@/utils/md5'

import { AGENT_TYPES, getAgentByAgentType } from '@/constants/platform/config'

import { channelApi } from '@/api/modules/channel'
import { agentApi } from '@/api'

defineProps({
  showChannelConfig: {
    type: Boolean,
    default: false,
  },
})

const store = useAgentFormStore()

const channelInfo = inject('channelConfig') || {}
const channelFormRef = ref()
const channelEditable = ref(false)
const channelForm = reactive({
  key: '',
  base_url: '',
  models: [],
  config: {
    agent_type: 'chat',
  },
})
const agentFormRef = ref()

const agentTypeOptions = [
  {
    icon: 'agent',
    label: window.$t('agent.dify.agent_type_chat'),
    description: window.$t('agent.dify.agent_type_chat_desc'),
    value: AGENT_TYPES.DIFY_AGENT,
  },
  {
    icon: 'completion-agent',
    label: window.$t('agent.dify.agent_type_workflow'),
    description: window.$t('agent.dify.agent_type_workflow_desc'),
    value: AGENT_TYPES.DIFY_WORKFLOW,
  },
]

const inputUpdateRequest = () => {
  return agentApi.dify.workflow_field_list(store.form_data.custom_config.channel_config.channel_id).then(res => {
    return res.user_input_form
      .map(item => {
        const type = Object.keys(item)[0]
        const value = Object.values(item)[0] as any
        if (!type) return null
        return {
          id: generateRandomId(6, true),
          variable: value.variable,
          type: type === 'paragraph' ? 'textarea' : type === 'select' ? 'select' : 'text',
          label: value.label,
          desc: value.desc,
          required: value.required,
          multiple: value.multiple || false,
          options: (value.options || []).map((item: string) => ({
            id: generateRandomId(6, true),
            label: item,
          })),
          max_length: value.max_length || 0,
          show_word_limit: value.show_word_limit || false,
          is_system: true,
        }
      })
      .filter(Boolean)
  })
}

const onChannelSave = async () => {
  const valid = await channelFormRef.value.validate()
  if (!valid) return
  const agent = getAgentByAgentType(store.agent_type)
  const model =
    (agent && agent.mode === 'completion' ? 'workflow-' : '') + md5(`${channelForm.key}_${channelForm.base_url}`)
  const name = 'dify'
  const saveData = {
    channel_id: channelInfo.value.channel_id,
    key: channelForm.key,
    base_url: channelForm.base_url,
    config: channelForm.config,
    models: [model],
    name,
  }
  const resultData = await channelApi.save({
    data: saveData,
  })
  Object.assign(channelInfo.value, resultData)
  if (!saveData.channel_id) saveData.channel_id = resultData.channel_id
  store.form_data.custom_config.channel_config = saveData
  store.form_data.model = model
  ElMessage.success(window.$t('action_save_success'))
  channelEditable.value = true
}

const validateForm = async () => {
  channelFormRef.value && channelFormRef.value.validate()
  if (agentFormRef.value) await agentFormRef.value.validate()
  return true
}

watch(
  () => store.agent_data,
  ({ channel_config = {} } = {}) => {
    channelEditable.value = !!+channel_config.channel_id
    channelInfo.value.channel_id = +channel_config.channel_id || 0
    channelInfo.value.key = channelForm.key = channel_config.key || ''
    channelInfo.value.base_url = channelForm.base_url = channel_config.base_url || 'https://api.dify.ai/v1'
    channelInfo.value.models = channelForm.models = channel_config.models || []
    channelInfo.value.config = channelForm.config = {
      ...(channel_config.config || {}),
      agent_type: channel_config.config?.agent_type || 'chat',
    }
  },
  { immediate: true, deep: true }
)

defineExpose({
  validateForm,
  onChannelSave,
})
</script>

<style scoped></style>
