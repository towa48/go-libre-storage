<template>
    <section class="body">
        <my-header></my-header>
        <div class="container-fluid">
            <div class="row">
                <my-sidebar></my-sidebar>

                <main role="main" class="col-md-9 col-lg-10">
                    <router-view></router-view>
                </main>
            </div>
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
</style>