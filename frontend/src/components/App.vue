<template>
  <div class="container">
    <h2>Offshoot Report</h2>

    <!-- 项目名称 -->
    <div class="form-group">
      <label for="projectTitle">请输入项目名称:</label>
      <input
          id="projectTitle"
          v-model="form.projectTitle"
          type="text"
          placeholder="项目名称"
      />
    </div>

    <!-- 备份数量 -->
    <div class="form-group">
      <label for="backups">请输入备份数量:</label>
      <input
          id="backups"
          v-model="form.backups"
          type="text"
          placeholder="备份数量"
      />
    </div>

    <!-- 选择MHL文件 -->
    <div class="form-group">
      <label>请选择Log文件:</label>
      <div class="file-select">
        <input
            type="text"
            v-model="form.filePath"
            readonly
            placeholder="未选择文件"
        />
        <button @click="selectFile">选择</button>
      </div>
    </div>

    <!-- 生成按钮 -->
    <button class="generate-btn" @click="generateReport">生成</button>

    <!-- 消息提示 -->
    <div v-if="message" :class="['message', messageClass]">
      {{ message }}
    </div>
  </div>
</template>

<script>
import { GenerateReport, SelectFiles } from "../../wailsjs/go/main/App"

export default {
  name: "App",
  data() {
    return {
      form: {
        projectTitle: "",
        backups: "",
        filePath: ""
      },
      message: "",
      messageClass: ""
    }
  },
  methods: {
    async selectFile() {
      try {
        const files = await SelectFiles()
        if (files && files.length > 0) {
          this.form.filePath = files[0]
        }
      } catch (error) {
        this.showMessage("选择文件时出错: " + error, "error")
      }
    },

    async generateReport() {
      if (!this.form.projectTitle || !this.form.backups || !this.form.filePath) {
        this.showMessage("请填写完整信息", "error")
        return
      }

      try {
        await GenerateReport({
          projectTitle: this.form.projectTitle,
          backups: this.form.backups,
          filePaths: [this.form.filePath]
        })
        this.showMessage("报告生成成功!" + this.form.projectTitle, "success")
      } catch (error) {
        this.showMessage("生成报告失败: " + error, "error")
      }
    },

    showMessage(text, type) {
      this.message = text
      this.messageClass = type
      setTimeout(() => {
        this.message = ""
        this.messageClass = ""
      }, 4000)
    }
  }
}
</script>

<style scoped>
.container {
  width: 400px;
  margin: 0 auto;
  padding: 10px;
  font-family: sans-serif;
}

h2 {
  text-align: center;
  margin-bottom: 15px;
}

.form-group {
  margin-bottom: 10px;
}

label {
  display: block;
  font-weight: bold;
  margin-bottom: 3px;
}

input[type="text"] {
  width: 100%;
  padding: 6px;
  font-size: 14px;
  border: 1px solid #ccc;
  border-radius: 4px;
  box-sizing: border-box;
}

.file-select {
  display: flex;
  gap: 8px;
}

.file-select input {
  flex: 1;
}

button {
  background: #007bff;
  color: white;
  border: none;
  padding: 6px 12px;
  border-radius: 4px;
  cursor: pointer;
  font-size: 14px;
}

button:hover {
  background: #0056b3;
}

.generate-btn {
  width: 100%;
  margin-top: 10px;
  background: #28a745;
  font-weight: bold;
}

.generate-btn:hover {
  background: #218838;
}

.message {
  margin-top: 10px;
  padding: 6px;
  border-radius: 4px;
  font-size: 14px;
  font-weight: bold;
}

.message.success {
  background: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.message.error {
  background: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}
</style>
