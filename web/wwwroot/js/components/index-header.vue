<template>
    <header class="header">
        <div class="logo-container">
            <a href="/" class="logo">GoLibreStorage</a>
            <!-- TODO: toggle sidebar here -->
        </div>
        <div class="header-right"></div>
    </header>
</template>

<script>
export default {
}
</script>

<style lang="scss">
.header {
    height: 60px;
    left: 0;
    position: fixed;
    right: 0;
    top: 0;
    z-index: 1020;

    background: $white;
    border-bottom: 1px solid #e9e9e6;
    border-top: 3px solid #ededed;

    .logo {
	    float: left;
	    margin: 10px 0 0 15px;

        img {
	        color: transparent;
        }
    }
}

.header-right {
    float: right;
    height: 56px;
}
</style>