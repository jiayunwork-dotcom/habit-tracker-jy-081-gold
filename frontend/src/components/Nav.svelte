<script>
  export let activeTab
</script>

<nav class="bg-white shadow-sm border-b border-gray-200 sticky top-0 z-10">
  <div class="container mx-auto px-4 max-w-4xl">
    <div class="flex items-center justify-between h-16">
      <div class="flex items-center space-x-2">
        <span class="text-2xl">🌟</span>
        <span class="text-xl font-bold text-gray-800">习惯养成</span>
      </div>
      
      <div class="flex space-x-1">
        <button 
          on:click={() => activeTab = 'habits'}
          class="px-4 py-2 rounded-lg font-medium transition-colors {activeTab === 'habits' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-600 hover:bg-gray-100'}"
        >
          习惯
        </button>
        <button 
          on:click={() => activeTab = 'stats'}
          class="px-4 py-2 rounded-lg font-medium transition-colors {activeTab === 'stats' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-600 hover:bg-gray-100'}"
        >
          统计
        </button>
        <button 
          on:click={() => activeTab = 'achievements'}
          class="px-4 py-2 rounded-lg font-medium transition-colors {activeTab === 'achievements' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-600 hover:bg-gray-100'}"
        >
          成就
        </button>
        <button 
          on:click={() => activeTab = 'settings'}
          class="px-4 py-2 rounded-lg font-medium transition-colors {activeTab === 'settings' ? 'bg-indigo-100 text-indigo-700' : 'text-gray-600 hover:bg-gray-100'}"
        >
          设置
        </button>
      </div>
    </div>
  </div>
</nav>
