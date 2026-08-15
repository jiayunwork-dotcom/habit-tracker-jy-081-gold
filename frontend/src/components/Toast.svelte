<script>
  export let toast
</script>

{#if toast}
  <div class="fixed bottom-4 right-4 z-50 animate-fade-in">
    <div class="px-6 py-3 rounded-lg shadow-lg {toast.type === 'success' ? 'bg-green-500' : 'bg-red-500'} text-white font-medium">
      {toast.message}
    </div>
  </div>
{/if}

<style>
  @keyframes fade-in {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
  .animate-fade-in {
    animation: fade-in 0.3s ease-out;
  }
</style>
