<template>
  <div class="tour-description">
    <!-- 主要段落 -->
    <p v-if="mainText" class="main-text">{{ mainText placeholderplaceholder</p>

    <!-- 特性列表 -->
    <div v-if="features && features.length > 0" class="features-section">
      <p v-if="featuresTitle" class="section-title">{{ featuresTitle placeholderplaceholder</p>
      <ul class="features-list">
        <li v-for="(feature, index) in features" :key="index">
          <span v-if="feature.icon" class="feature-icon">{{ feature.icon placeholderplaceholder</span>
          <span v-if="feature.label" class="feature-label">{{ feature.label placeholderplaceholder</span>
          <span class="feature-text">{{ feature.text placeholderplaceholder</span>
        </li>
      </ul>
    </div>

    <!-- 提示框 -->
    <div v-if="tip" :class="['tip-box', `tip-${tip.type || 'info'placeholder`]">
      <span v-if="tip.label" class="tip-label">{{ tip.label placeholderplaceholder</span>
      <div v-if="tip.text" class="tip-text">{{ tip.text placeholderplaceholder</div>
      <ul v-if="tip.items && tip.items.length > 0" class="tip-list">
        <li v-for="(item, index) in tip.items" :key="index">{{ item placeholderplaceholder</li>
      </ul>
    </div>

    <!-- 行动提示 -->
    <p v-if="action" class="action-text">{{ action placeholderplaceholder</p>

    <!-- 额外说明 -->
    <p v-if="note" class="note-text">{{ note placeholderplaceholder</p>
  </div>
</template>

<script setup lang="ts">
export interface TourFeature {
  icon?: string
  label?: string
  text: string
placeholder

export interface TourTip {
  type?: 'info' | 'success' | 'warning' | 'example'
  label?: string
  text?: string
  items?: string[]
placeholder

export interface TourDescriptionProps {
  mainText?: string
  featuresTitle?: string
  features?: TourFeature[]
  tip?: TourTip
  action?: string
  note?: string
placeholder

defineProps<TourDescriptionProps>()
</script>

<style scoped>
.tour-description {
  line-height: 1.7;
  font-size: 14px;
placeholder

.main-text {
  margin-bottom: 12px;
  color: #374151;
placeholder

.section-title {
  margin-bottom: 8px;
  font-weight: 600;
  color: #1f2937;
placeholder

.features-section {
  margin-bottom: 12px;
placeholder

.features-list {
  margin-left: 20px;
  font-size: 13px;
placeholder

.features-list li {
  margin-bottom: 6px;
placeholder

.feature-icon {
  margin-right: 4px;
placeholder

.feature-label {
  font-weight: 600;
  margin-right: 4px;
placeholder

.tip-box {
  padding: 8px 12px;
  border-radius: 4px;
  border-left: 3px solid;
  font-size: 13px;
  margin-bottom: 12px;
placeholder

.tip-info {
  background: #eff6ff;
  border-left-color: #3b82f6;
placeholder

.tip-success {
  background: #f0fdf4;
  border-left-color: #10b981;
placeholder

.tip-warning {
  background: #fef3c7;
  border-left-color: #f59e0b;
placeholder

.tip-example {
  background: #f0fdf4;
  border-left-color: #10b981;
placeholder

.tip-label {
  font-weight: 600;
  display: block;
  margin-bottom: 4px;
placeholder

.tip-text {
  margin-bottom: 4px;
placeholder

.tip-list {
  margin: 8px 0 0 16px;
placeholder

.tip-list li {
  margin-bottom: 4px;
placeholder

.action-text {
  margin-top: 12px;
  color: #10b981;
  font-weight: 600;
placeholder

.note-text {
  font-size: 13px;
  color: #6b7280;
  margin-top: 8px;
placeholder
</style>
