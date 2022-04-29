package cluster

import (
	"kubespace/server/global"
	"kubespace/server/model/kubernetes"
)

func CreateK8SCluster(cluster kubernetes.K8SCluster) (err error) {
	err = global.GVA_DB.Create(&cluster).Error
	return
}

func ListK8SCluster(p *kubernetes.PaginationQ, k *[]kubernetes.K8SCluster) (err error) {

	if p.Page < 1 {
		p.Page = 1
	}
	if p.Size < 1 {
		p.Size = 10
	}

	offset := p.Size * (p.Page - 1)
	tx := global.GVA_DB
	if p.Keyword != "" {
		tx = global.GVA_DB.Where("cluster_name like ?", "%"+p.Keyword+"%").Limit(p.Size).Offset(offset).Find(&k)
	} else {
		tx = global.GVA_DB.Limit(p.Size).Offset(offset).Find(&k)

	}

	var total int64
	tx.Count(&total)
	//p.Total = tx.RowsAffected
	p.Total = total
	if err := tx.Error; err != nil {
		return err
	}

	return nil
}

func GetK8sCluster(id uint) (K8sCluster kubernetes.K8SCluster, err error) {
	err = global.GVA_DB.Where("id = ?", id).First(&K8sCluster).Error
	if err != nil {
		return K8sCluster, err
	}
	return K8sCluster, nil
}

func DelCluster(ids kubernetes.ClusterIds) (err error) {
	var k kubernetes.K8SCluster
	err2 := global.GVA_DB.Delete(&k, ids.Data)
	if err2.Error != nil {
		return err2.Error
	}
	return nil
}
