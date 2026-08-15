<script>
  import { createEventDispatcher } from 'svelte'
  const dispatch = createEventDispatcher()

  export let habits = []
  let draggedIndex = null

  function handleDragStart(index) {
    draggedIndex = index
  }

  function handleDragOver(e, index) {
    e.preventDefault()
    if (draggedIndex !== null && draggedIndex !== index) {
      const newHabits = [...habits]
      const dragged = newHabits.splice(draggedIndex, 1)[0]
      newHabits.splice(index, 0, dragged)
      habits = newHabits
      draggedIndex = index
    }
  }

  function handleDragEnd() {
    dispatch('reorder', habits.map(h => h.id))
    draggedIndex = null
  }
</script>

<div class="space-y-3">
  {#each habits as habit, index (habit.id)}
    <div 
      class="bg-white rounded-xl shadow-sm border border-gray-200 p-4 hover:shadow-md transition-shadow cursor-move"
      draggable="true"
      on:dragstart={() => handleDragStart(index)}
      on:dragover={(e) => handleDragOver(e, index)}
      on:dragend={handleDragEnd}
    >
      <div class="flex items-center justify-between">
        <div class="flex items-center space-x-3">
          <div 
            class="w-10 h-10 rounded-full flex items-center justify-center text-white font-bold"
            style="background-color: {habit.color}"
          >
            {habit.name.charAt(0)}
          </div>
          <div>
            <h3 class="font-semibold text-gray-800">{habit.name}</h3>
            <div class="text-sm text-gray-500">
              {#if habit.current_streak > 0}
                <span class="text-orange-500 font-medium">🔥 连续 {habit.current_streak} 天</span>
              {:else}
                <span>开始你的第一天吧！</span>
              {/if}
              {#if habit.longest_streak > 0 && habit.longest_streak > habit.current_streak}
                <span class="ml-2">· 最长 {habit.longest_streak} 天</span>
              {/if}
            </div>
          </div>
        </div>
        
        <div class="flex items-center space-x-2">
          {#if habit.today_checked_in}
            <button 
              on:click={() => dispatch('cancel', { id: habit.id })}
              class="px-4 py-2 bg-green-100 text-green-700 rounded-lg font-medium hover:bg-green-200 transition-colors"
            >
              ✓ 已完成
            </button>
          {:else}
            <button 
              on:click={() => dispatch('checkin', { id: habit.id })}
              class="px-4 py-2 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition-colors"
            >
              打卡
            </button>
          {/if}
          
          <div class="relative group">
            <button class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg">
              ⋮
            </button>
            <div class="absolute right-0 mt-2 w-32 bg-white rounded-lg shadow-lg border border-gray-200 opacity-0 invisible group-hover:opacity-100 group-hover:visible transition-all z-10">
              <button 
                on:click={() => dispatch('edit', habit)}
                class="w-full px-4 py-2 text-left text-gray-700 hover:bg-gray-50 rounded-t-lg"
              >
                编辑
              </button>
              <button 
                on:click={() => dispatch('delete', habit.id)}
                class="w-full px-4 py-2 text-left text-red-600 hover:bg-red-50 rounded-b-lg"
              >
                删除
              </button>
            </div>
          </div>
        </div>
      </div>
    </div>
  {/each}
</div>
