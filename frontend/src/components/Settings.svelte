<script>
  import { onMount, createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  let doNotDisturbEnabled = false
  let doNotDisturbStart = '23:00'
  let doNotDisturbEnd = '07:00'
  let saving = false

  async function fetchSettings() {
    const res = await fetch('/api/settings')
    const data = await res.json()
    if (data.do_not_disturb_enabled !== undefined) {
      doNotDisturbEnabled = data.do_not_disturb_enabled === 'true'
    }
    if (data.do_not_disturb_start) {
      doNotDisturbStart = data.do_not_disturb_start
    }
    if (data.do_not_disturb_end) {
      doNotDisturbEnd = data.do_not_disturb_end
    }
  }

  async function saveSettings() {
    saving = true
    await fetch('/api/settings', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        do_not_disturb_enabled: doNotDisturbEnabled.toString(),
        do_not_disturb_start: doNotDisturbStart,
        do_not_disturb_end: doNotDisturbEnd
      })
    })
    dispatch('saved')
    saving = false
  }

  onMount(() => {
    fetchSettings()
  })
</script>

<div class="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
  <h2 class="text-xl font-bold text-gray-800 mb-6">设置</h2>
  
  <div class="space-y-6">
    <div class="border-b border-gray-200 pb-6">
      <h3 class="text-lg font-semibold text-gray-700 mb-4">免打扰模式</h3>
      
      <div class="flex items-center justify-between mb-4">
        <div>
          <p class="font-medium text-gray-700">启用免打扰</p>
          <p class="text-sm text-gray-500">在指定时间内不会收到提醒通知</p>
        </div>
        <button
          on:click={() => doNotDisturbEnabled = !doNotDisturbEnabled}
          class="relative inline-flex h-6 w-11 items-center rounded-full transition-colors {doNotDisturbEnabled ? 'bg-indigo-600' : 'bg-gray-300'}"
        >
          <span
            class="inline-block h-4 w-4 transform rounded-full bg-white transition-transform {doNotDisturbEnabled ? 'translate-x-6' : 'translate-x-1'}"
          />
        </button>
      </div>

      {#if doNotDisturbEnabled}
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">开始时间</label>
            <input
              type="time"
              bind:value={doNotDisturbStart}
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-1">结束时间</label>
            <input
              type="time"
              bind:value={doNotDisturbEnd}
              class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
            />
          </div>
        </div>
      {/if}
    </div>

    <div class="flex justify-end">
      <button
        on:click={saveSettings}
        disabled={saving}
        class="px-4 py-2 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors disabled:opacity-50"
      >
        {saving ? '保存中...' : '保存设置'}
      </button>
    </div>
  </div>
</div>
