<template>  
<header class="navbar bg-base-100 shadow-xl fixed">
    <nav class="navbar-start gap-2">   
        <slot name="start"></slot>
    </nav>
    <nav class="navbar-center gap-2">  
        <slot name="center"></slot>
    </nav>
    <nav class="navbar-end gap-2">    
        <slot name="end"></slot>
    </nav>
</header>
</template>
