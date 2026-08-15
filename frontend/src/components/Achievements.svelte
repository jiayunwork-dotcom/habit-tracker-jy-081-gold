<script>
  import { onMount } from 'svelte'

  let achievements = []

  async function fetchAchievements() {
    const res = await fetch('/api/stats/achievements')
    achievements = await res.json()
  }

  function getProgressColor(achievement) {
    if (achievement.earned) return 'bg-green-500'
    const progress = (achievement.progress / achievement.target) * 100
    if (progress < 30) return 'bg-gray-300'
    if (progress < 70) return 'bg-yellow-500'
    return 'bg-green-400'
  }

  $: earnedCount = achievements.filter(a => a.earned).length
  $: totalCount = achievements.length

  onMount(() => {
    fetchAchievements()
  })
</script>

<div>
  <h1 class="text-2xl font-bold text-gray-800 mb-6">成就系统</h1>

  <div class="bg-gradient-to-r from-indigo-500 to-purple-600 p-6 rounded-xl text-white mb-8">
    <div class="flex items-center justify-between">
      <div>
        <div class="text-4xl font-bold">{earnedCount}/{totalCount}</div>
        <div class="opacity-80">已解锁成就</div>
      </div>
      <div class="text-6xl">🏆</div>
    </div>
    <div class="mt-4 bg-white/20 rounded-full h-2">
      <div 
        class="bg-white rounded-full h-2 transition-all duration-500"
        style="width: {totalCount > 0 ? (earnedCount / totalCount) * 100 : 0}%"
      />
    </div>
  </div>

  <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
    {#each achievements as achievement}
      <div class="bg-white rounded-xl shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow">
        <div class="flex items-start space-x-4">
          <div class="text-4xl {achievement.earned ? '' : 'grayscale opacity-50'}">
            {achievement.icon}
          </div>
          <div class="flex-1">
            <div class="flex items-center justify-between">
              <h3 class="font-semibold text-gray-800 {achievement.earned ? '' : 'opacity-50'}">
                {achievement.name}
              </h3>
              {#if achievement.earned}
                <span class="px-2 py-1 bg-green-100 text-green-700 text-xs font-medium rounded-full">
                  已获得
                </span>
              {/if}
            </div>
            <p class="text-sm text-gray-500 mt-1">{achievement.description}</p>
            
            {#if !achievement.earned && achievement.target > 1}
              <div class="mt-3">
                <div class="flex justify-between text-xs text-gray-500 mb-1">
                  <span>进度</span>
                  <span>{achievement.progress}/{achievement.target}</span>
                </div>
                <div class="bg-gray-200 rounded-full h-2">
                  <div 
                    class="{getProgressColor(achievement)} rounded-full h-2 transition-all duration-500"
                    style="width: {Math.min((achievement.progress / achievement.target) * 100, 100)}%"
                  />
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>
    {/each}
  </div>
</div>
