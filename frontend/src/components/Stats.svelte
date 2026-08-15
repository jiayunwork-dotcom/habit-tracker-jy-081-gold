<script>
  import { onMount } from 'svelte'
  import { format, subDays } from 'date-fns'

  export let habits = []
  let heatmapData = {}
  let overview = {}

  $: maxStreak = habits.reduce((max, h) => Math.max(max, h.longest_streak || 0), 0)

  async function fetchStats() {
    const [heatmapRes, overviewRes] = await Promise.all([
      fetch('/api/stats/heatmap'),
      fetch('/api/stats/overview')
    ])
    heatmapData = await heatmapRes.json()
    overview = await overviewRes.json()
  }

  function getHeatColor(percent) {
    if (!percent) return 'bg-gray-200'
    if (percent < 33) return 'bg-green-200'
    if (percent < 66) return 'bg-green-400'
    return 'bg-green-600'
  }

  function getDays() {
    const days = []
    for (let i = 364; i >= 0; i--) {
      const date = subDays(new Date(), i)
      days.push({
        date: format(date, 'yyyy-MM-dd'),
        display: format(date, 'MM/dd')
      })
    }
    return days
  }

  onMount(() => {
    fetchStats()
  })
</script>

<div>
  <h1 class="text-2xl font-bold text-gray-800 mb-6">统计概览</h1>

  <div class="grid grid-cols-3 gap-4 mb-8">
    <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
      <div class="text-3xl font-bold text-indigo-600">{habits.length || 0}</div>
      <div class="text-gray-500">活跃习惯</div>
    </div>
    <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
      <div class="text-3xl font-bold text-green-600">{overview.completed_today || 0}</div>
      <div class="text-gray-500">今日完成</div>
    </div>
    <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
      <div class="text-3xl font-bold text-orange-600">{maxStreak}</div>
      <div class="text-gray-500">最高连续</div>
    </div>
  </div>

  <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200 mb-8">
    <h2 class="text-lg font-semibold text-gray-800 mb-4">打卡热力图</h2>
    <div class="overflow-x-auto">
      <div class="flex flex-wrap gap-1" style="width: 750px;">
        {#each getDays() as day}
          <div 
            class="w-3 h-3 rounded-sm {getHeatColor(heatmapData[day.date])}"
            title="{day.display}: {heatmapData[day.date] || 0}%"
          />
        {/each}
      </div>
    </div>
    <div class="flex items-center gap-2 mt-4 text-sm text-gray-500">
      <span>少</span>
      <div class="w-3 h-3 rounded-sm bg-gray-200"></div>
      <div class="w-3 h-3 rounded-sm bg-green-200"></div>
      <div class="w-3 h-3 rounded-sm bg-green-400"></div>
      <div class="w-3 h-3 rounded-sm bg-green-600"></div>
      <span>多</span>
    </div>
  </div>

  <div class="bg-white p-6 rounded-xl shadow-sm border border-gray-200">
    <h2 class="text-lg font-semibold text-gray-800 mb-4">习惯详情</h2>
    <div class="space-y-4">
      {#each habits as habit}
        <div class="flex items-center justify-between p-4 bg-gray-50 rounded-lg">
          <div class="flex items-center space-x-3">
            <div 
              class="w-8 h-8 rounded-full flex items-center justify-center text-white font-bold text-sm"
              style="background-color: {habit.color}"
            >
              {habit.name.charAt(0)}
            </div>
            <div>
              <div class="font-medium text-gray-800">{habit.name}</div>
              <div class="text-sm text-gray-500">
                累计 {habit.total_checkins} 次 · 连续 {habit.current_streak} 天
              </div>
            </div>
          </div>
          <div class="text-right">
            <div class="text-lg font-bold text-gray-800">{habit.longest_streak}</div>
            <div class="text-xs text-gray-500">最长连续</div>
          </div>
        </div>
      {/each}
    </div>
  </div>
</div>
