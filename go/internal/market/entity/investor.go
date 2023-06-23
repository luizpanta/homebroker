package entity

// Estrutura dados do Investidor ID, Nome e Posição de ativos
type Investor struct {
	ID            string
	Name          string
	AssetPosition []*InvestorAssetPosition // Slice (go)
}

// Criar novo investidor ID e posição em branco
func NewInvestor(id string) *Investor {
	return &Investor{
		ID:            id,
		AssetPosition: []*InvestorAssetPosition{}, // Array dinâmico
	}
}

// Método adicionar posição
func (i *Investor) AddAssetPosition(assetPosition *InvestorAssetPosition) {
	i.AssetPosition = append(i.AssetPosition, assetPosition) //Adiciona em []*InvestorAssetPosition{}
}

// Método que atualiza as posições
func (i *Investor) UpdateAssetPosition(assetID string, qtdShares int) {
	assetPosition := i.GetAssetPosition(assetID)
	if assetPosition == nil {
		i.AssetPosition = append(i.AssetPosition, NewInvestorAssetPosition(assetID, qtdShares)) // Adiciona uma primeira vez
	} else {
		assetPosition.Shares += qtdShares //Soma a posição existente
	}
}

// Método que busca a posição atual
func (i *Investor) GetAssetPosition(assetID string) *InvestorAssetPosition {
	for _, assetPosition := range i.AssetPosition {
		if assetPosition.AssetID == assetID {
			return assetPosition
		}
	}
	return nil
}

// Estrutura dados para posição de ativos do investidor
type InvestorAssetPosition struct {
	AssetID string
	Shares  int
}

// Função que cria nova posição
func NewInvestorAssetPosition(assetID string, shares int) *InvestorAssetPosition {
	return &InvestorAssetPosition{
		AssetID: assetID,
		Shares:  shares,
	}
}
