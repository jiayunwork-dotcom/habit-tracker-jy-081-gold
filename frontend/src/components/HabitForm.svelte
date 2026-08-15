<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  export let habit = null

  let name = habit?.name || ''
  let color = habit?.color || '#4F46E5'
  let frequencyType = habit?.frequency_type || 'daily'
  let frequencyValue = habit?.frequency_value || 1
  let specificDays = habit?.specific_days || ''
  let reminderTime = habit?.reminder_time || ''
  let notificationPermission = 'default'

  if ('Notification' in window) {
    notificationPermission = Notification.permission
  }

  async function requestNotificationPermission() {
    if (!('Notification' in window)) {
      alert('您的浏览器不支持通知功能')
      return
    }
    const permission = await Notification.requestPermission()
    notificationPermission = permission
  }

  const colors = [
    '#4F46E5', '#10B981', '#F59E0B', '#EF4444', '#8B5CF6', 
    '#EC4899', '#06B6D4', '#84CC16', '#F97316', '#6366F1'
  ]

  const frequencyOptions = [
    { value: 'daily', label: '每天打卡' },
    { value: 'weekly_n', label: '每周N次' },
    { value: 'specific_days', label: '指定星期' },
    { value: 'monthly_n', label: '每月N次' }
  ]

  const weekDays = [
    { value: '1', label: '周一' },
    { value: '2', label: '周二' },
    { value: '3', label: '周三' },
    { value: '4', label: '周四' },
    { value: '5', label: '周五' },
    { value: '6', label: '周六' },
    { value: '7', label: '周日' }
  ]

  function toggleDay(day) {
    const days = specificDays ? specificDays.split(',') : []
    const index = days.indexOf(day)
    if (index > -1) {
      days.splice(index, 1)
    } else {
      days.push(day)
    }
    specificDays = days.sort().join(',')
  }

  function handleSubmit(e) {
    e.preventDefault()
    dispatch('submit', {
      name,
      color,
      frequency_type: frequencyType,
      frequency_value: parseInt(frequencyValue),
      specific_days: frequencyType === 'specific_days' ? specificDays : '',
      reminder_time: reminderTime || null
    })
  }
</script>

<form on:submit={handleSubmit} class="space-y-4">
  <div>
    <label class="block text-sm font-medium text-gray-700 mb-1">习惯名称</label>
    <input 
      type="text" 
      bind:value={name}
      required
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
      placeholder="例如：每天阅读30分钟"
    />
  </div>

  <div>
    <label class="block text-sm font-medium text-gray-700 mb-2">颜色</label>
    <div class="flex flex-wrap gap-2">
      {#each colors as c}
        <button 
          type="button"
          on:click={() => color = c}
          class="w-8 h-8 rounded-full border-2 transition-transform hover:scale-110 {color === c ? 'border-gray-800 scale-110' : 'border-transparent'}"
          style="background-color: {c}"
        />
      {/each}
    </div>
  </div>

  <div>
    <label class="block text-sm font-medium text-gray-700 mb-1">打卡频率</label>
    <select 
      bind:value={frequencyType}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
    >
      {#each frequencyOptions as opt}
        <option value={opt.value}>{opt.label}</option>
      {/each}
    </select>
  </div>

  {#if frequencyType === 'weekly_n' || frequencyType === 'monthly_n'}
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-1">
        {frequencyType === 'weekly_n' ? '每周次数' : '每月次数'}
      </label>
      <input 
        type="number" 
        bind:value={frequencyValue}
        min="1"
        max={frequencyType === 'weekly_n' ? 7 : 31}
        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500"
      />
    </div>
  {/if}

  {#if frequencyType === 'specific_days'}
    <div>
      <label class="block text-sm font-medium text-gray-700 mb-2">选择星期</label>
      <div class="flex flex-wrap gap-2">
        {#each weekDays as day}
          <button 
            type="button"
            on:click={() => toggleDay(day.value)}
            class="px-4 py-2 rounded-lg font-medium transition-colors {specificDays.includes(day.value) ? 'bg-indigo-600 text-white' : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}"
          >
            {day.label}
          </button>
        {/each}
      </div>
    </div>
  {/if}

  <div>
    <label class="block text-sm font-medium text-gray-700 mb-1">提醒时间（可选）</label>
    <input 
      type="time" 
      bind:value={reminderTime}
      class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:ring-2 focus:ring-indigo-500 focus:border-indigo-500"
    />
    {#if reminderTime && notificationPermission !== 'granted'}
      <div class="mt-2 p-3 bg-amber-50 border border-amber-200 rounded-lg">
        <p class="text-sm text-amber-700 mb-2">
          ⚠️ 需要开启浏览器通知权限才能收到提醒
        </p>
        <button
          on:click={requestNotificationPermission}
          class="text-sm font-medium text-amber-700 hover:text-amber-800 underline"
        >
          点击开启通知权限
        </button>
      </div>
    {/if}
  </div>

  <div class="flex gap-3 pt-4">
    <button 
      type="button"
      on:click={() => dispatch('cancel')}
      class="flex-1 px-4 py-2 border border-gray-300 rounded-lg font-medium text-gray-700 hover:bg-gray-50 transition-colors"
    >
      取消
    </button>
    <button 
      type="submit"
      class="flex-1 px-4 py-2 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors"
    >
      {habit ? '更新' : '创建'}习惯
    </button>
  </div>
</form>
