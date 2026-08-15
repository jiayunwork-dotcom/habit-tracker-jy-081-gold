<script>
  import { onMount, onDestroy } from 'svelte'
  import HabitList from './components/HabitList.svelte'
  import HabitForm from './components/HabitForm.svelte'
  import Stats from './components/Stats.svelte'
  import Achievements from './components/Achievements.svelte'
  import Settings from './components/Settings.svelte'
  import Nav from './components/Nav.svelte'
  import Toast from './components/Toast.svelte'

  let habits = []
  let activeTab = 'habits'
  let showForm = false
  let editingHabit = null
  let toast = null
  let notificationPermission = 'default'
  let settings = {}
  let reminderCheckInterval = null

  const encouragements = [
    '加油！你可以的！💪',
    '坚持就是胜利！🌟',
    '好习惯正在养成！🌱',
    '今天也要努力哦！✨',
    '你正在变得更好！💫'
  ]

  function getRandomEncouragement() {
    return encouragements[Math.floor(Math.random() * encouragements.length)]
  }

  async function requestNotificationPermission() {
    if (!('Notification' in window)) {
      showToast('您的浏览器不支持通知功能', 'warning')
      return
    }
    const permission = await Notification.requestPermission()
    notificationPermission = permission
    if (permission === 'granted') {
      showToast('通知权限已开启！')
    } else if (permission === 'denied') {
      showToast('通知权限被拒绝，请在浏览器设置中手动开启', 'warning')
    }
  }

  function isInDoNotDisturb() {
    if (!settings.do_not_disturb_enabled || settings.do_not_disturb_enabled !== 'true') {
      return false
    }

    const now = new Date()
    const currentMinutes = now.getHours() * 60 + now.getMinutes()
    
    const [startHour, startMin] = (settings.do_not_disturb_start || '23:00').split(':').map(Number)
    const [endHour, endMin] = (settings.do_not_disturb_end || '07:00').split(':').map(Number)
    
    const startMinutes = startHour * 60 + startMin
    const endMinutes = endHour * 60 + endMin

    if (startMinutes <= endMinutes) {
      return currentMinutes >= startMinutes && currentMinutes < endMinutes
    } else {
      return currentMinutes >= startMinutes || currentMinutes < endMinutes
    }
  }

  function checkReminders() {
    if (notificationPermission !== 'granted') return
    if (isInDoNotDisturb()) return

    const now = new Date()
    const currentTime = `${String(now.getHours()).padStart(2, '0')}:${String(now.getMinutes()).padStart(2, '0')}`

    habits.forEach(habit => {
      if (!habit.reminder_time || habit.today_checked_in) return
      if (habit.reminder_time === currentTime && !habit.is_archived) {
        new Notification('习惯提醒', {
          body: `${habit.name} - ${getRandomEncouragement()}`,
          icon: 'data:image/svg+xml,<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">🌟</text></svg>'
        })
      }
    })
  }

  async function fetchHabits() {
    const res = await fetch('/api/habits')
    habits = await res.json()
  }

  async function fetchSettings() {
    const res = await fetch('/api/settings')
    settings = await res.json()
  }

  function showToast(message, type = 'success') {
    toast = { message, type }
    setTimeout(() => toast = null, 3000)
  }

  async function createHabit(habit) {
    const res = await fetch('/api/habits', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(habit)
    })
    if (res.ok) {
      showToast('习惯创建成功！')
      fetchHabits()
      showForm = false
    }
  }

  async function updateHabit(id, habit) {
    const res = await fetch(`/api/habits/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(habit)
    })
    if (res.ok) {
      showToast('习惯更新成功！')
      fetchHabits()
      editingHabit = null
      showForm = false
    }
  }

  async function deleteHabit(id) {
    if (confirm('确定要删除这个习惯吗？所有打卡记录也会被删除。')) {
      const res = await fetch(`/api/habits/${id}`, { method: 'DELETE' })
      if (res.ok) {
        showToast('习惯已删除')
        fetchHabits()
      }
    }
  }

  async function checkIn(id, date = null) {
    const url = date ? `/api/habits/${id}/checkin?date=${date}` : `/api/habits/${id}/checkin`
    const res = await fetch(url, { method: 'POST' })
    const data = await res.json()
    if (res.ok) {
      showToast(data.encouragement)
      fetchHabits()
      if (data.new_achievements && data.new_achievements.length > 0) {
        setTimeout(() => {
          showToast(`🎉 获得新成就：${data.new_achievements[0].achievement.name}！`)
        }, 1500)
      }
    }
  }

  async function cancelCheckIn(id, date = null) {
    const url = date ? `/api/habits/${id}/checkin?date=${date}` : `/api/habits/${id}/checkin`
    const res = await fetch(url, { method: 'DELETE' })
    if (res.ok) {
      showToast('已取消打卡')
      fetchHabits()
    }
  }

  async function reorderHabits(ids) {
    await fetch('/api/habits/reorder', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ ids })
    })
    fetchHabits()
  }

  onMount(() => {
    if ('Notification' in window) {
      notificationPermission = Notification.permission
    }
    fetchHabits()
    fetchSettings()
    
    reminderCheckInterval = setInterval(() => {
      fetchHabits()
      fetchSettings()
      checkReminders()
    }, 60000)
  })

  onDestroy(() => {
    if (reminderCheckInterval) {
      clearInterval(reminderCheckInterval)
    }
  })
</script>

<div class="min-h-screen bg-gray-50">
  <Nav bind:activeTab />
  
  <main class="container mx-auto px-4 py-8 max-w-4xl">
    {#if activeTab === 'habits'}
      <div class="mb-6 flex justify-between items-center">
        <h1 class="text-2xl font-bold text-gray-800">我的习惯</h1>
        <button 
          on:click={() => { showForm = true; editingHabit = null }}
          class="bg-indigo-600 hover:bg-indigo-700 text-white px-4 py-2 rounded-lg font-medium transition-colors"
        >
          + 添加习惯
        </button>
      </div>

      {#if showForm}
        <div class="mb-6 bg-white p-6 rounded-xl shadow-sm border border-gray-200">
          <h2 class="text-lg font-semibold mb-4">{editingHabit ? '编辑习惯' : '创建新习惯'}</h2>
          <HabitForm 
            habit={editingHabit}
            on:submit={(e) => {
              if (editingHabit) {
              updateHabit(editingHabit.id, e.detail)
            } else {
              createHabit(e.detail)
            }
          }}
          on:cancel={() => { showForm = false; editingHabit = null }}
          />
        </div>
      {/if}

      {#if habits.length === 0}
        <div class="text-center py-12 bg-white rounded-xl shadow-sm p-8">
        <div class="text-6xl mb-4">🌱</div>
        <h3 class="text-xl font-semibold text-gray-700 mb-2">还没有习惯</h3>
        <p class="text-gray-500 mb-4">点击上方按钮创建你的第一个习惯吧！</p>
        </div>
      {:else}
        <HabitList 
          {habits}
          on:checkin={(e) => checkIn(e.detail.id, e.detail.date)}
          on:cancel={(e) => cancelCheckIn(e.detail.id, e.detail.date)}
          on:edit={(e) => { editingHabit = e.detail; showForm = true }}
          on:delete={(e) => deleteHabit(e.detail)}
          on:reorder={(e) => reorderHabits(e.detail)}
        />
      {/if}
    {:else if activeTab === 'stats'}
      <Stats {habits} />
    {:else if activeTab === 'achievements'}
      <Achievements />
    {:else if activeTab === 'settings'}
      <Settings on:saved={() => { showToast('设置已保存！'); fetchSettings() }} />
    {/if}
  </main>

  <Toast bind:toast />
</div>
