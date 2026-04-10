<template>
  <div class="app-container">
    <van-nav-bar :title="pageTitle" fixed />

    <div class="content">
      <van-steps :active="currentPhase" active-color="#1989fa">
        <van-step>选择选题</van-step>
        <van-step>添加选项</van-step>
        <van-step>评分</van-step>
        <van-step>结果</van-step>
      </van-steps>

      <div v-if="currentPhase === 0" class="step-content">
        <h3>选择或创建选题</h3>
        
        <div class="action-buttons">
          <van-field v-model="newTopic" placeholder="输入新选题名称" @keyup.enter="addTopic">
            <template #button>
              <van-button size="small" type="primary" @click="addTopic">添加选题</van-button>
            </template>
          </van-field>
          <van-button type="default" block @click="openTemplatePicker">
            从模板选择
          </van-button>
        </div>

        <div class="topics-list">
          <van-cell-group>
            <van-cell 
              v-for="topic in topics" 
              :key="topic.id" 
              clickable
              @click="toggleTopic(topic.id)"
            >
              <template #title>
                <div class="topic-title-cell">
                  <span>{{ topic.title }}</span>
                  <van-icon name="edit" class="topic-edit-icon" @click.stop="openEditTopicModal(topic)" />
                </div>
              </template>
              <template #label>
                <span class="topic-is-template" v-if="topic.isTemplate">模板</span>
              </template>
              <template #value>
                <div class="cell-actions">
                  <van-checkbox 
                    :model-value="selectedTopicIds.includes(topic.id)" 
                    shape="square"
                    @click.stop="toggleTopic(topic.id)"
                  />
                  <van-icon 
                    v-if="!topic.isTemplate"
                    name="award" 
                    class="save-template-icon" 
                    title="存为模板"
                    @click.stop="saveAsTemplate(topic.id)"
                  />
                  <van-icon 
                    name="delete" 
                    class="delete-icon" 
                    @click.stop="deleteTopic(topic.id)"
                  />
                </div>
              </template>
            </van-cell>
          </van-cell-group>
        </div>

        <van-button type="primary" block @click="goToAddOptions" :disabled="selectedTopicIds.length === 0">
          下一步
        </van-button>
      </div>

      <div v-if="currentPhase === 1" class="step-content">
        <div class="topic-tabs">
          <div 
            v-for="(topicId, idx) in selectedTopicIds" 
            :key="topicId" 
            class="topic-tab"
            :class="{ active: currentTopicIndex === idx }"
            @click="onTopicTabChange(idx)"
          >
            {{ getTopicTitle(topicId) }}
          </div>
        </div>

        <h3>为 {{ getTopicTitle(selectedTopicIds[currentTopicIndex]) }} 添加选项</h3>
        
        <div class="topic-options-section">
          <div class="topic-header">
            <van-button size="small" type="default" @click="showAddOptionModal">
              + 单个添加
            </van-button>
          </div>
          
          <van-field
            v-model="multiOptionsInput"
            type="textarea"
            placeholder="每行一个选项，可一次输入多个"
            :rows="3"
            autosize
            @keyup.enter="addMultiOptions"
          />
          <van-button size="small" type="primary" @click="addMultiOptions" style="margin-top: 8px;">
            批量添加
          </van-button>
          
          <div class="options-list">
            <div v-for="(opt, idx) in getTopicCurrentRound().options" :key="idx" class="option-item">
              <span>{{ opt.title }}</span>
              <van-icon name="close" @click="removeOption(idx)" />
            </div>
          </div>
          
          <div v-if="!getTopicCurrentRound().options || getTopicCurrentRound().options.length === 0" class="no-options">
            请添加至少2个选项
          </div>
        </div>

        <van-button type="primary" block @click="goToScoring" :disabled="!currentTopicHasOptions">
          下一步
        </van-button>
        
        <van-button type="default" block @click="currentPhase = 0" style="margin-top: 8px;">
          上一步
        </van-button>
      </div>

      <div v-if="currentPhase === 2" class="step-content">
<div class="topic-tabs">
          <div 
            v-for="(topicId, idx) in selectedTopicIds" 
            :key="topicId" 
            class="topic-tab"
            :class="{ active: currentTopicIndex === idx }"
            @click="onTopicTabChange(idx)"
          >
            {{ getTopicTitle(topicId) }}
            <span v-if="getTopicRoundStatus(topicId) === 'completed'" class="round-check">✓</span>
          </div>
        </div>

        <h3>评分 - {{ getTopicTitle(selectedTopicIds[currentTopicIndex]) }} 第{{ currentRoundIndex + 1 }}轮</h3>
        
        <div class="round-indicator">
          <div class="round-tabs">
            <div 
              v-for="(round, idx) in getTopicRounds(selectedTopicIds[currentTopicIndex])" 
              :key="idx" 
              class="round-tab"
              :class="{ active: currentRoundIndex === idx }"
              @click="currentRoundIndex = idx"
            >
              <span>第{{ idx + 1 }}轮</span>
              <span v-if="round.status === 'completed'" class="round-check">✓</span>
              <van-icon 
                v-if="round.status !== 'completed' && round.status !== 'locked' && idx > 0" 
                name="close" 
                size="12" 
                class="round-delete" 
                @click.stop="deleteRound(idx)" 
              />
            </div>
            <div class="round-tab add-round" @click="addRound" v-if="canAddRound">
              + 新一轮
            </div>
          </div>
        </div>
        
        <p class="hint">重要性(1-5分) + 必要性(0或1)</p>

        <div class="topic-options-section">
          <div v-for="(opt, idx) in getTopicCurrentRound().options" :key="idx" class="score-item">
            <div class="option-name">{{ opt.title }}</div>
            <div class="score-controls">
              <van-rate v-model="opt.importance" :max="5" size="18" />
              <van-radio-group v-model="opt.necessity" direction="horizontal">
                <van-radio :name="1" shape="square">必要</van-radio>
                <van-radio :name="0" shape="square">不必要</van-radio>
              </van-radio-group>
            </div>
          </div>
        </div>

        <van-button type="primary" block @click="calculateResults">
          计算结果
        </van-button>
        
        <van-button type="default" block @click="currentPhase = 1" style="margin-top: 8px;">
          上一步
        </van-button>
      </div>

      <div v-if="currentPhase === 3" class="step-content">
<div class="topic-tabs">
          <div 
            v-for="(topicId, idx) in selectedTopicIds" 
            :key="topicId" 
            class="topic-tab"
            :class="{ active: currentTopicIndex === idx }"
            @click="onTopicTabChange(idx)"
          >
            {{ getTopicTitle(topicId) }}
            <span v-if="idx === currentTopicIndex && getTopicRoundStatus() === 'completed'" class="round-check">✓</span>
          </div>
        </div>

        <div class="round-indicator">
          <div class="round-tabs">
            <div 
              v-for="(round, idx) in getTopicRounds(selectedTopicIds[currentTopicIndex])" 
              :key="idx" 
              class="round-tab"
              :class="{ active: currentRoundIndex === idx }"
              @click="currentRoundIndex = idx"
            >
              第{{ idx + 1 }}轮
              <span v-if="round.status === 'completed'" class="round-check">✓</span>
            </div>
          </div>
        </div>

        <h3>筛选结果 - {{ getTopicTitle(selectedTopicIds[currentTopicIndex]) }} 第{{ currentRoundIndex + 1 }}轮</h3>

        <div class="topic-result-section">
          <van-card v-for="(r, idx) in getTopicCurrentRound().results" :key="idx" class="result-card">
            <template #title>
              <div class="result-title">
                <span>{{ r.option }}</span>
                <van-tag :type="r.isWinner ? 'success' : 'default'">
                  {{ r.isWinner ? '最佳' : '第' + r.rank + '名' }}
                </van-tag>
              </div>
            </template>
            <template #desc>
              <div class="result-detail">
                <div>重要性: {{ r.importance }}分</div>
                <div>必要性: {{ r.necessity === 1 ? '必要' : '不必要' }}</div>
                <div class="total-score">总分: {{ r.totalScore }}</div>
              </div>
            </template>
          </van-card>

          <div class="topic-match-score">
            <span>匹配度: </span>
            <van-progress :percentage="getTopicCurrentRound().matchScore" :stroke-width="8" />
            <span class="match-percent">{{ getTopicCurrentRound().matchScore }}%</span>
          </div>
        </div>

        <div class="overall-match">
          <h4>整体匹配度: {{ overallMatchScore }}%</h4>
        </div>

        <van-button type="primary" block @click="finishCurrentRound" v-if="currentRoundStatus !== 'completed'">
          完成本轮
        </van-button>

        <van-button type="default" block @click="goBackToScoring" style="margin-top: 8px;">
          重新评分
        </van-button>

        <van-button type="warning" block @click="reset" style="margin-top: 8px;">
          重新开始
        </van-button>
      </div>
    </div>

    <van-popup v-model:show="showOptionModal" position="bottom" round>
      <div class="option-modal">
        <h3>添加选项 - {{ getTopicTitle(selectedTopicIds[currentTopicIndex]) }}</h3>
        <van-field v-model="newOption" placeholder="输入选项名称" @keyup.enter="confirmAddOption" />
        <van-button type="primary" block @click="confirmAddOption">确定</van-button>
      </div>
    </van-popup>

    <van-popup v-model:show="showTemplatePicker" position="bottom" round>
      <div class="template-modal">
        <h3>选择模板选题</h3>
        <div class="template-list">
          <van-checkbox-group v-model="selectedTemplateIds">
            <van-cell-group>
              <van-cell 
                v-for="tpl in templateTopics" 
                :key="tpl.id" 
                clickable 
                @click="toggleTemplate(tpl.id)"
              >
                <template #title>
                  <div>{{ tpl.title }}</div>
                  <div class="template-desc">{{ tpl.description }}</div>
                </template>
                <template #label>
                  <div v-if="selectedTemplateIds.includes(tpl.id)" class="template-edit">
                    <van-field v-model="templateNames[tpl.id]" placeholder="自定义选题名称(可选)" size="small" @click.stop />
                  </div>
                </template>
                <template #value>
                  <div class="template-actions">
                    <van-icon name="delete" class="template-delete" @click.stop="deleteTemplate(tpl.id)" />
                    <van-checkbox :name="tpl.id" shape="square" />
                  </div>
                </template>
              </van-cell>
            </van-cell-group>
          </van-checkbox-group>
        </div>
        <van-button type="primary" block @click="importSelectedTemplates" :disabled="selectedTemplateIds.length === 0">
          导入选中选题 ({{ selectedTemplateIds.length }})
        </van-button>
      </div>
    </van-popup>

    <van-popup v-model:show="showSaveTemplateModal" position="bottom" round>
      <div class="option-modal">
        <h3>存为模板</h3>
        <van-field v-model="templateTitleInput" placeholder="输入模板名称" />
        <van-button type="primary" block @click="confirmSaveAsTemplate">确定</van-button>
      </div>
    </van-popup>

    <van-popup v-model:show="showEditTopicModal" position="bottom" round>
      <div class="option-modal">
        <h3>编辑选题</h3>
        <van-field v-model="editTopicTitle" placeholder="输入选题名称" />
        <van-button type="primary" block @click="confirmEditTopic">确定</van-button>
      </div>
    </van-popup>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

const API_BASE = 'http://localhost:8888/api'

const currentPhase = ref(0)
const currentTopicIndex = ref(0)
// 新增：定义 currentRoundIndex 响应式变量（核心修复）
const currentRoundIndex = ref(0)
const topics = ref([])
const selectedTopicIds = ref([])
const newTopic = ref('')
const newOption = ref('')
const multiOptionsInput = ref('')
const topicRounds = ref({})
const topicRoundIndexes = ref({})  // 每个topic独立的轮数索引
const showOptionModal = ref(false)
const showTemplatePicker = ref(false)
const showSaveTemplateModal = ref(false)
const templateTitleInput = ref('')
const savingTemplateTopicId = ref(null)
const showEditTopicModal = ref(false)
const editTopicTitle = ref('')
const editingTopicId = ref(null)

function openEditTopicModal(topic) {
  editingTopicId.value = topic.id
  editTopicTitle.value = topic.title
  showEditTopicModal.value = true
}

const pageTitle = computed(() => {
  const titles = ['', '选择选题', '添加选项', '评分', '结果']
  return titles[currentPhase.value] || '多轮打分筛选'
})

function getTopicRounds(topicId) {
  return topicRounds.value[topicId] || []
}

function getTopicCurrentRound(topicId) {
  if (!topicId) {
    if (selectedTopicIds.value.length === 0 || currentTopicIndex.value >= selectedTopicIds.value.length) {
      return { options: [], results: [], matchScore: 0, status: 'pending' }
    }
    topicId = selectedTopicIds.value[currentTopicIndex.value]
  }
  const rounds = topicRounds.value[topicId]
  if (!rounds || rounds.length === 0) {
    return { options: [], results: [], matchScore: 0, status: 'pending' }
  }
  //const roundIdx = currentRoundIndex.value
  // 新增：校验 currentRoundIndex 边界，避免越界
  const roundIdx = Math.max(0, Math.min(currentRoundIndex.value, rounds.length - 1))
  const round = rounds[roundIdx]
  if (!round) {
    return { options: [], results: [], matchScore: 0, status: 'pending' }
  }
  return round
}

function getTopicRoundStatus(topicId) {
  if (!topicId) {
    if (selectedTopicIds.value.length === 0 || currentTopicIndex.value >= selectedTopicIds.value.length) {
      return 'pending'
    }
    topicId = selectedTopicIds.value[currentTopicIndex.value]
  }
  const rounds = topicRounds.value[topicId]
  if (!rounds || rounds.length === 0) return 'pending'
  const roundIdx = currentRoundIndex.value
  return rounds[roundIdx]?.status || 'pending'
}

// 切换 topic tab 时加载数据
async function onTopicTabChange(idx) {
  const previousTopicId = selectedTopicIds.value[currentTopicIndex.value]
  if (previousTopicId) {
    topicRoundIndexes.value[previousTopicId] = currentRoundIndex.value
    await saveOptions(previousTopicId)
  }
  
  currentTopicIndex.value = idx
  const topicId = selectedTopicIds.value[idx]
  if (topicRoundIndexes.value[topicId] === undefined) {
    topicRoundIndexes.value[topicId] = 0
  }
  await loadTopicData(topicId)
}

const currentRoundStatus = computed(() => {
  if (selectedTopicIds.value.length === 0) return 'pending'
  return getTopicRoundStatus(selectedTopicIds.value[currentTopicIndex.value])
})

const canAddRound = computed(() => {
  if (selectedTopicIds.value.length === 0) return false
  const topicId = selectedTopicIds.value[currentTopicIndex.value]
  const rounds = topicRounds.value[topicId]
  if (!rounds || rounds.length === 0) return false
  return currentRoundIndex.value === rounds.length - 1
})

async function loadTopics() {
  try {
    console.log('Loading topics...')
    const res = await fetch(`${API_BASE}/topics`)
    console.log('Topics response status:', res.status)
    const data = await res.json()
    console.log('Topics data:', data)
    if (data.code === 0 && data.data) {
      topics.value = data.data.filter(t => !t.isTemplate)
    }
  } catch (e) {
    console.error('Failed to load topics', e)
  }
}

async function loadTemplateTopics() {
  try {
    const res = await fetch(`${API_BASE}/topics/templates`)
    const data = await res.json()
    if (data.code === 0 && data.data) {
      templateTopics.value = data.data
    }
  } catch (e) {
    console.error('Failed to load templates', e)
  }
}

async function openTemplatePicker() {
  if (templateTopics.value.length === 0) {
    await loadTemplateTopics()
  }
  showTemplatePicker.value = true
}

function toggleTemplate(id) {
  const idx = selectedTemplateIds.value.indexOf(id)
  if (idx > -1) {
    selectedTemplateIds.value.splice(idx, 1)
  } else {
    selectedTemplateIds.value.push(id)
  }
}

async function importSelectedTemplates() {
  if (selectedTemplateIds.value.length === 0) return
  try {
    const importData = selectedTemplateIds.value.map(id => ({
      templateId: id,
      title: templateNames.value[id] || ''
    }))
    const res = await fetch(`${API_BASE}/topics/import`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(importData)
    })
    const data = await res.json()
    if (data.code === 0 && data.data) {
      for (const t of data.data) {
        topics.value.push(t)
        selectedTopicIds.value.push(t.id)
        topicRounds.value[t.id] = [{
          options: [],
          results: [],
          matchScore: 0,
          importanceStatus: 'pending',
          necessityStatus: 'pending',
          status: 'pending'
        }]
      }
    }
    showTemplatePicker.value = false
    selectedTemplateIds.value = []
    templateNames.value = {}
  } catch (e) {
    console.error('Failed to import templates', e)
    alert('导入失败: ' + e.message)
  }
}

async function deleteTemplate(id) {
  if (!confirm('确定要删除这个模板吗？')) return
  try {
    const res = await fetch(`${API_BASE}/topics/${id}`, { method: 'DELETE' })
    const data = await res.json()
    if (data.code === 0) {
      const idx = templateTopics.value.findIndex(t => t.id === id)
      if (idx > -1) templateTopics.value.splice(idx, 1)
      const selIdx = selectedTemplateIds.value.indexOf(id)
      if (selIdx > -1) selectedTemplateIds.value.splice(selIdx, 1)
    }
  } catch (e) {
    console.error('Failed to delete template', e)
    alert('删除失败')
  }
}

async function toggleTopic(id) {
  const idx = selectedTopicIds.value.indexOf(id)
  if (idx > -1) {
    selectedTopicIds.value.splice(idx, 1)
    delete topicRounds.value[id]
  } else {
    selectedTopicIds.value.push(id)
    topicRounds.value[id] = [{
      options: [],
      results: [],
      matchScore: 0,
      importanceStatus: 'pending',
      necessityStatus: 'pending',
      status: 'pending'
    }]
    await loadTopicData(id)
  }
}

async function deleteTopic(id) {
  if (!confirm('确定要删除这个选题吗？')) return
  try {
    const res = await fetch(`${API_BASE}/topics/${id}`, { method: 'DELETE' })
    const data = await res.json()
    if (data.code === 0) {
      const idx = topics.value.findIndex(t => t.id === id)
      if (idx > -1) topics.value.splice(idx, 1)
      const selIdx = selectedTopicIds.value.indexOf(id)
      if (selIdx > -1) {
        selectedTopicIds.value.splice(selIdx, 1)
        delete topicRounds.value[id]
      }
    }
  } catch (e) {
    console.error('Failed to delete topic', e)
    alert('删除失败')
  }
}

async function saveAsTemplate(id) {
  savingTemplateTopicId.value = id
  const t = topics.value.find(tp => tp.id === id)
  templateTitleInput.value = t ? t.title : ''
  showSaveTemplateModal.value = true
}

async function confirmSaveAsTemplate() {
  if (!templateTitleInput.value.trim()) {
    alert('请输入模板名称')
    return
  }
  try {
    const res = await fetch(`${API_BASE}/topics/${savingTemplateTopicId.value}/save-as-template`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: templateTitleInput.value.trim() })
    })
    const data = await res.json()
    if (data.code === 0) {
      alert('已存为模板')
    }
  } catch (e) {
    console.error('Failed to save as template', e)
    alert('保存失败')
  }
  showSaveTemplateModal.value = false
  templateTitleInput.value = ''
  savingTemplateTopicId.value = null
}

async function confirmEditTopic() {
  if (!editTopicTitle.value.trim()) {
    alert('请输入选题名称')
    return
  }
  try {
    const res = await fetch(`${API_BASE}/topics/${editingTopicId.value}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ title: editTopicTitle.value.trim() })
    })
    const data = await res.json()
    if (data.code === 0) {
      const idx = topics.value.findIndex(t => t.id === editingTopicId.value)
      if (idx > -1) {
        topics.value[idx].title = editTopicTitle.value.trim()
      }
    }
  } catch (e) {
    console.error('Failed to update topic', e)
    alert('保存失败')
  }
  showEditTopicModal.value = false
  editTopicTitle.value = ''
  editingTopicId.value = null
}

async function loadTopicData(topicId) {
  try {
    console.log('Loading topic data for:', topicId)
    const [optsRes, roundsRes] = await Promise.all([
      fetch(`${API_BASE}/topics/${topicId}/options`),
      fetch(`${API_BASE}/topics/${topicId}/rounds`)
    ])
    const optsData = await optsRes.json()
    const roundsData = await roundsRes.json()
    console.log('Options data:', optsData)
    console.log('Rounds data:', roundsData)
    
    let options = []
    if (optsData.code === 0 && optsData.data) {
      options = optsData.data.map(o => ({
        title: o.title,
        importance: o.importance || 0,
        necessity: o.necessity || 0
      }))
    }
    console.log('Base options:', options)
    
    if (roundsData.code === 0 && roundsData.data && roundsData.data.length > 0) {
      console.log('Loading', roundsData.data.length, 'rounds for topic', topicId)
      topicRounds.value[topicId] = roundsData.data.map((r, rIdx) => {
        let results = []
        try {
          results = r.results ? JSON.parse(r.results) : []
        } catch (e) {
          results = []
        }
        
        let roundOptions = []
        if (r.options) {
          try {
            const parsedOptions = typeof r.options === 'string' ? JSON.parse(r.options) : r.options
            if (Array.isArray(parsedOptions) && parsedOptions.length > 0) {
              roundOptions = parsedOptions.map(o => ({
                title: o.title || '',
                importance: o.importance || 0,
                necessity: o.necessity || 0
              }))
            }
          } catch (e) {
            roundOptions = []
          }
        }
        
        if (roundOptions.length === 0 && options.length > 0) {
          roundOptions = options.map(o => ({ ...o }))
        }
        
        return {
          id: r.id,
          options: roundOptions,
          results: results,
          matchScore: r.matchScore || 0,
          status: r.status || 'pending',
          roundNumber: r.roundNumber
        }
      })
    } else if (options.length > 0) {
      topicRounds.value[topicId] = [{
        options: options.map(o => ({ ...o })),
        results: [],
        matchScore: 0,
        importanceStatus: 'pending',
        necessityStatus: 'pending',
        status: 'pending'
      }]
    }
  } catch (e) {
    console.error('Failed to load topic data', e)
  }
}

function getTopicTitle(id) {
  const t = topics.value.find(t => t.id === id)
  return t ? t.title : ''
}

async function addTopic() {
  if (newTopic.value.trim()) {
    try {
      console.log('Adding topic:', newTopic.value)
      const res = await fetch(`${API_BASE}/topics`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ title: newTopic.value.trim() })
      })
      console.log('Response status:', res.status)
      const data = await res.json()
      console.log('Response data:', data)
      if (data.code === 0 && data.data) {
        topics.value.push(data.data)
        selectedTopicIds.value.push(data.data.id)
        topicRounds.value[data.data.id] = [{
          options: [],
          results: [],
          matchScore: 0,
          importanceStatus: 'pending',
          necessityStatus: 'pending',
          status: 'pending'
        }]
      }
    } catch (e) {
      console.error('Failed to add topic', e)
      alert('添加失败: ' + e.message)
    }
    newTopic.value = ''
  }
}

function showAddOptionModal() {
  newOption.value = ''
  showOptionModal.value = true
}

async function addMultiOptions() {
  if (!multiOptionsInput.value.trim() || selectedTopicIds.value.length === 0) return
  
  const topicId = selectedTopicIds.value[currentTopicIndex.value]
  const currentRound = topicRounds.value[topicId][currentRoundIndex.value]
  
  const lines = multiOptionsInput.value.split('\n').map(s => s.trim()).filter(s => s)
  const newOptions = []
  
  for (const title of lines) {
    if (title) {
      const opt = {
        title: title,
        importance: 0,
        necessity: 0
      }
      currentRound.options.push(opt)
      newOptions.push(opt)
      syncNewOptionToAllRounds(topicId, opt)
    }
  }
  
  multiOptionsInput.value = ''
  if (newOptions.length > 0) {
    await saveOptions(topicId)
  }
}

// 同步新选项到所有未完成的轮次
function syncNewOptionToAllRounds(topicId, newOption) {
  const rounds = topicRounds.value[topicId]
  if (!rounds) return
  for (let i = 0; i < rounds.length; i++) {
    // 跳过已完成和locked的轮次
    if (rounds[i].status === 'completed' || rounds[i].status === 'locked') continue
    // 检查是否已存在同名选项
    const exists = rounds[i].options.some(o => o.title === newOption.title)
    if (!exists) {
      rounds[i].options.push({ ...newOption })
    }
  }
}

async function confirmAddOption() {
  if (newOption.value.trim() && selectedTopicIds.value.length > 0) {
    const topicId = selectedTopicIds.value[currentTopicIndex.value]
    const currentRound = topicRounds.value[topicId][currentRoundIndex.value]
    console.log('Adding option, topicId:', topicId, 'round:', currentRound)
    const opt = {
      title: newOption.value.trim(),
      importance: 0,
      necessity: 0
    }
    currentRound.options.push(opt)
    syncNewOptionToAllRounds(topicId, opt)
    console.log('Options after push:', currentRound.options)
    newOption.value = ''
    showOptionModal.value = false
    await saveOptions(topicId)
    console.log('Saved, current options:', currentRound.options)
  }
}

async function removeOption(idx) {
  const topicId = selectedTopicIds.value[currentTopicIndex.value]
  const rounds = topicRounds.value[topicId]
  const removedTitle = rounds[currentRoundIndex.value].options[idx].title
  
  rounds[currentRoundIndex.value].options.splice(idx, 1)
  
  for (let i = 0; i < rounds.length; i++) {
    if (rounds[i].status === 'completed' || rounds[i].status === 'locked') continue
    const optIdx = rounds[i].options.findIndex(o => o.title === removedTitle)
    if (optIdx > -1) {
      rounds[i].options.splice(optIdx, 1)
    }
  }
  
  await saveOptions(topicId)
}

async function saveOptions(topicId) {
  const round = topicRounds.value[topicId][currentRoundIndex.value]
  console.log('saveOptions round:', round)
  console.log('saveOptions topicId:', topicId, 'round index:', currentRoundIndex.value)
  const options = round.options.map(o => ({
    title: o.title,
    importance: o.importance,
    necessity: o.necessity
  }))
  console.log('Saving options:', options)
  try {
    const res = await fetch(`${API_BASE}/topics/${topicId}/options`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(options)
    })
    const data = await res.json()
    console.log('Save options response:', data)
  } catch (e) {
    console.error('Failed to save options', e)
  }
}

const currentTopicHasOptions = computed(() => {
  if (selectedTopicIds.value.length === 0) return false
  const round = getTopicCurrentRound(selectedTopicIds.value[currentTopicIndex.value])
  return round.options && round.options.length >= 2
})

async function goToAddOptions() {
  // 进入添加选项界面时重新加载选项数据（强制刷新）
  for (const topicId of selectedTopicIds.value) {
    await loadTopicData(topicId)
  }
  currentTopicIndex.value = 0
  currentPhase.value = 1
}

async function goToScoring() {
  for (const topicId of selectedTopicIds.value) {
    await saveOptions(topicId)
  }
  await saveAllRounds()
  for (const topicId of selectedTopicIds.value) {
    await loadTopicData(topicId)
  }
  topicRoundIndexes.value = {}
  currentRoundIndex.value = 0
  currentPhase.value = 2
}

function goBackToScoring() {
  for (const topicId of selectedTopicIds.value) {
    topicRounds.value[topicId][currentRoundIndex.value].status = 'pending'
  }
  currentPhase.value = 2
}

async function addRound() {
  const topicId = selectedTopicIds.value[currentTopicIndex.value]
  
  const res = await fetch(`${API_BASE}/topics/${topicId}/options`)
  const data = await res.json()
  
  let allOptions = []
  if (data.code === 0 && data.data) {
    allOptions = data.data.map(o => ({
      title: o.title,
      importance: 0,
      necessity: 0
    }))
  }
  
  const newRound = {
    options: allOptions,
    results: [],
    matchScore: 0,
    importanceStatus: 'pending',
    necessityStatus: 'pending',
    status: 'pending'
  }
  topicRounds.value[topicId].push(newRound)
  currentRoundIndex.value = topicRounds.value[topicId].length - 1
  await saveAllRounds()
}

async function deleteRound(idx) {
  const topicId = selectedTopicIds.value[currentTopicIndex.value]
  const rounds = topicRounds.value[topicId]
  
  if (idx === 0) {
    alert('第一轮不能删除')
    return
  }
  
  if (rounds[idx].status === 'completed' || rounds[idx].status === 'locked') {
    alert('已完成的轮次不能删除')
    return
  }
  
  try {
    const roundId = rounds[idx].id
    if (roundId) {
      await fetch(`${API_BASE}/rounds/${roundId}`, { method: 'DELETE' })
    }
  } catch (e) {
    console.error('Failed to delete round from backend', e)
  }
  
  rounds.splice(idx, 1)
  if (currentRoundIndex.value >= rounds.length) {
    currentRoundIndex.value = rounds.length - 1
  }
  await saveAllRounds()
}

async function calculateResults() {
  for (const topicId of selectedTopicIds.value) {
    const round = topicRounds.value[topicId][currentRoundIndex.value]
    
    if (round.options.some(o => o.importance === 0)) {
      alert('请完成所有重要性评分')
      return
    }
    if (round.options.some(o => o.necessity === undefined || o.necessity === null)) {
      alert('请完成所有必要性评分')
      return
    }

    const resultsData = round.options.map(opt => ({
      option: opt.title,
      importance: opt.importance,
      necessity: opt.necessity,
      totalScore: opt.importance * opt.necessity
    }))

    resultsData.sort((a, b) => b.totalScore - a.totalScore)

    round.results = resultsData.map((r, idx) => ({
      ...r,
      rank: idx + 1,
      isWinner: idx === 0 && r.totalScore > 0
    }))

    round.matchScore = calculateMatchScore(resultsData)
  }
  currentPhase.value = 3
  for (const topicId of selectedTopicIds.value) {
    await saveOptions(topicId)
  }
  await saveAllRounds()
}

function calculateMatchScore(results) {
  if (results.length === 0) return 0
  const totalWeight = results.reduce((sum, r) => sum + r.importance, 0)
  if (totalWeight === 0) return 0
  const necessaryWeight = results
    .filter(r => r.necessity === 1)
    .reduce((sum, r) => sum + r.importance, 0)
  return Math.round((necessaryWeight / totalWeight) * 100)
}

const overallMatchScore = computed(() => {
  if (selectedTopicIds.value.length === 0) return 0
  const scores = selectedTopicIds.value.map(tid => getTopicCurrentRound(tid).matchScore)
  const sum = scores.reduce((a, b) => a + b, 0)
  return Math.round(sum / scores.length)
})

function finishCurrentRound() {
  for (const topicId of selectedTopicIds.value) {
    topicRounds.value[topicId][currentRoundIndex.value].status = 'completed'
  }
  saveAllRounds()
}

async function saveAllRounds() {
  for (const topicId of selectedTopicIds.value) {
    const rounds = topicRounds.value[topicId].map((r, idx) => {
      const roundData = {
        roundNumber: idx + 1,
        importanceStatus: r.importanceStatus,
        necessityStatus: r.necessityStatus,
        status: r.status,
        matchScore: r.matchScore,
        options: r.options.map(o => ({
          title: o.title,
          importance: o.importance,
          necessity: o.necessity
        })),
        results: r.results
      }
      return roundData
    })
    try {
      await fetch(`${API_BASE}/topics/${topicId}/rounds`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(rounds)
      })
    } catch (e) {
      console.error('Failed to save rounds', e)
    }
  }
}

function reset() {
  currentPhase.value = 0
  // 重置所有 topic 的轮数索引
  topicRoundIndexes.value = {}
  currentTopicIndex.value = 0
  selectedTopicIds.value = []
  topicRounds.value = {}
}

onMounted(() => {
  loadTopics()
})
</script>

<style>
body { margin: 0; padding: 0; }
.app-container { min-height: 100vh; background: #f5f5f5; }
.content { padding: 16px; padding-top: 80px; }
.step-content { margin-top: 20px; }
.hint { color: #666; font-size: 14px; margin-bottom: 16px; }

.topics-list { margin-bottom: 16px; }

.topic-is-template { font-size: 12px; color: #969799; }
.topic-title-cell { display: flex; align-items: center; gap: 8px; }
.topic-edit-icon { color: #969799; cursor: pointer; margin-left: 4px; }
.topic-edit-icon:hover { color: #1989fa; }

.cell-actions { display: flex; align-items: center; gap: 8px; }
.save-template-icon { color: #969799; cursor: pointer; padding: 4px; }
.save-template-icon:hover { color: #ff976a; }
.delete-icon { color: #969799; cursor: pointer; padding: 4px; }
.delete-icon:hover { color: #ee0a24; }

.topic-tabs { display: flex; gap: 8px; margin-bottom: 16px; flex-wrap: wrap; }
.topic-tab {
  padding: 8px 16px; background: #fff; border-radius: 16px; font-size: 14px; cursor: pointer;
  display: flex; align-items: center; gap: 4px;
}
.topic-tab.active { background: #1989fa; color: #fff; }
.round-check { color: #07c160; }

.round-indicator { margin-bottom: 16px; }
.round-tabs { display: flex; gap: 8px; flex-wrap: wrap; }
.round-tab { 
  padding: 6px 12px; background: #fff; border-radius: 12px; font-size: 12px; cursor: pointer;
  display: flex; align-items: center; gap: 4px;
}
.round-tab.active { background: #1989fa; color: #fff; }
.round-tab.add-round { border: 1px dashed #1989fa; color: #1989fa; }
.round-tab { position: relative; }
.round-delete { 
  margin-left: 4px; color: #999; cursor: pointer; 
}
.round-delete:hover { color: #f44; }

.topic-options-section {
  background: #fff; border-radius: 8px; padding: 16px; margin-bottom: 16px;
}
.topic-header {
  display: flex; justify-content: space-between; align-items: center; margin-bottom: 12px;
}
.topic-title { font-weight: bold; font-size: 16px; }

.options-list { display: flex; flex-wrap: wrap; gap: 8px; margin-top: 8px; }
.option-item {
  display: flex; align-items: center; gap: 8px;
  padding: 8px 12px; background: #f5f5f5; border-radius: 16px;
}
.no-options { color: #999; font-size: 14px; text-align: center; padding: 16px; }

.score-item {
  display: flex; flex-direction: column;
  padding: 12px 0; border-bottom: 1px solid #eee;
}
.score-item:last-child { border-bottom: none; }
.option-name { font-weight: bold; margin-bottom: 8px; }
.score-controls { display: flex; justify-content: space-between; align-items: center; }

.topic-result-section { margin-bottom: 24px; }
.topic-result-header { margin-bottom: 12px; }

.result-card { margin-bottom: 8px; }
.result-title { display: flex; justify-content: space-between; align-items: center; }
.result-detail { color: #666; font-size: 14px; }
.total-score { font-weight: bold; color: #1989fa; margin-top: 4px; }

.topic-match-score {
  display: flex; align-items: center; gap: 12px; margin-top: 12px; padding: 12px;
  background: #f0f0f0; border-radius: 8px;
}
.match-percent { font-weight: bold; color: #1989fa; }

.overall-match {
  background: #fff; padding: 16px; border-radius: 8px; margin: 16px 0;
  text-align: center;
}
.overall-match h4 { margin: 0; color: #1989fa; }

.option-modal { padding: 16px; }
.option-modal h3 { margin: 0 0 16px 0; }

.action-buttons { margin-bottom: 16px; }

.template-modal { padding: 16px; max-height: 60vh; }
.template-modal h3 { margin: 0 0 16px 0; }
.template-list { max-height: 50vh; overflow-y: auto; margin-bottom: 16px; }
.template-desc { font-size: 12px; color: #999; margin-top: 4px; }
.template-edit { margin-top: 8px; }
.template-actions { display: flex; align-items: center; gap: 8px; }
.template-delete { color: #969799; cursor: pointer; padding: 4px; }
.template-delete:hover { color: #ee0a24; }
</style>