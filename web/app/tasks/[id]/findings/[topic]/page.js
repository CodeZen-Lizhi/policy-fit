export default function FindingDetailPage() {
  return (
    <main style={{ maxWidth: 960, margin: "24px auto", padding: 16 }}>
      <h1>🔴 高血压风险</h1>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 12 }}>
        <h2>风险说明</h2>
        <p>体检报告显示血压偏高且有复查建议，与保单既往症定义可能存在冲突。</p>
      </section>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 12 }}>
        <h2>体检证据</h2>
        <p>[para_12] 血压 155/95 mmHg，建议复查</p>
      </section>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 12 }}>
        <h2>条款证据</h2>
        <p>[para_120] 既往症是指在合同生效前已存在的疾病或异常</p>
      </section>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16, marginBottom: 12 }}>
        <h2>召回片段来源</h2>
        <ul>
          <li>source: policy_v20260228.pdf / section: 第1条既往症定义 / score: 0.91</li>
          <li>source: policy_v20260228.pdf / section: 第4条责任免除 / score: 0.86</li>
        </ul>
      </section>
      <section style={{ background: "#fff", borderRadius: 12, padding: 16 }}>
        <h2>追问清单</h2>
        <ul>
          <li>是否已被医生明确诊断为高血压？</li>
          <li>是否正在长期服用降压药物？</li>
          <li>体检日期是否在本次投保前？</li>
        </ul>
      </section>
    </main>
  );
}
