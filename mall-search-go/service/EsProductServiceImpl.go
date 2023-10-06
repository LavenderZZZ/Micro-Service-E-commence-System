package service

import "mall-search-go/store"

type EsProductServiceImpl struct {
	prouductDao   store.EsproductDao
	elasticClient *elastic.ElasticClient
}

func (s EsProductServiceImpl) ImportAll() (int, error) {
	esProductList, err := s.prouductDao.GetAllProductList(nil)
	if err != nil {
		return 0, err
	}

}
