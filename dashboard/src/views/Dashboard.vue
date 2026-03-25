<template>
  <AdminLayout 
    :baseUrl="apiUrl" 
    @logout="logout" 
    @refresh="fetchComments(pagination.page)"
  >
    <div v-if="loading" class="flex justify-center py-20">
      <div class="animate-spin rounded-full h-10 w-10 border-4 border-blue-500 border-t-transparent"></div>
    </div>
    
    <CommentList 
      v-else
      :data="tableData" 
      :pagination="pagination"
      @page-change="handlePageChange"
      @update="updateStatus"
    />
  </AdminLayout>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { ElMessage } from 'element-plus';
import request from '../utils/request';
import AdminLayout from '../components/AdminLayout.vue';
import CommentList from '../components/CommentList.vue';

const router = useRouter();
const loading = ref(false);
const tableData = ref([]);
const pagination = ref({ page: 1, limit: 10, totalPage: 1 });
const apiUrl = ref(localStorage.getItem('apiUrl') || window.location.origin);

const fetchComments = async (page = 1) => {
  loading.value = true;
  try {
    const res = await request.get(`/admin/comments/list?page=${page}`);
    if (res.data) {
      tableData.value = res.data.comments;
      pagination.value = res.data.pagination;
    }
  } catch (error) {
    ElMessage.error('加载失败');
  } finally {
    loading.value = false;
  }
};

const updateStatus = async (id, status) => {
  try {
    await request.put(`/admin/comments/status?id=${id}&status=${status}`);
    ElMessage.success('操作成功');
    fetchComments(pagination.value.page);
  } catch (error) {
    ElMessage.error('更新失败');
  }
};

const handlePageChange = (p) => fetchComments(p);

const logout = () => {
  localStorage.removeItem('token');
  router.push('/login');
};

onMounted(() => fetchComments(1));
</script>