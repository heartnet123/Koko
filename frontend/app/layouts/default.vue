<template>
  <div class="min-h-screen bg-default text-default font-sans flex selection:bg-primary-100 selection:text-primary-900">
    <AppSidebar />
    <main class="flex-1 flex flex-col min-w-0">
      <AppHeader />
      <div class="flex-1 overflow-y-auto pb-12">
        <slot />
      </div>
    </main>
  </div>
</template>
