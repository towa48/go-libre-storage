<template>
    <section class="body">
        <my-header></my-header>
        <div class="inner-wrapper">
            <my-sidebar></my-sidebar>

            <main role="main" class="content-body">
                <router-view></router-view>
            </main>
        </div>
    </section>
</template>

<script>
import HeaderComponent from './index-header.vue';
import SidebarComponent from './index-sidebar.vue';

export default {
    data: function() {
        return {}
    },
    components: {
        'my-header': HeaderComponent,
        'my-sidebar': SidebarComponent
    }
}
</script>

<style lang="scss">
.body {
    min-height: 100vh;
    width: 100%;
}
.inner-wrapper {
    min-height: 100vh;
    width: 100%;
    overflow: hidden;
    padding-top: 60px;
}
.content-body {
    display: table-cell;
    padding: 40px;
    position: relative;
    vertical-align: top;
}
@include media-breakpoint-up(md) {
    .content-body {
        display: block;
        margin-left: 300px;
    }
}
</style>