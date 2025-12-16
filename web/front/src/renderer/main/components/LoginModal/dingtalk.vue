<template>
  <div :style="{ height, width }" class="flex flex-col justify-center items-center">
    <div class="w-[220px] h-[220px] border relative rounded-lg overflow-hidden flex items-center justify-center bg-gray-50">
      <svg-icon name="dingtalk" size="64" />
    </div>
    <p class="text-sm text-regular mt-3 text-opacity-60">{{ $t('login.login_by_dingtalk') }}</p>
    <div class="mt-4">
        <ElButton type="primary" @click="handleLogin"> {{ $t('login.immediate_login') }} </ElButton>
    </div>
  </div>
</template>

<script lang="ts" setup>
import SvgIcon from '@/components/SvgIcon.vue'

const props = withDefaults(
  defineProps<{
    width?: string
    height?: string
  }>(),
  {
    width: '100%',
    height: '280px'
  }
)

const handleLogin = () => {
  // TODO: Implement DingTalk login logic
  console.log('DingTalk login')
}
</script>
