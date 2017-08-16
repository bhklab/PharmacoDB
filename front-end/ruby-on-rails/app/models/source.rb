class Source < ApplicationRecord
	has_many :source_cell_names
	has_many :source_tissue_names
	has_many :source_drug_names
	has_many :cells, through: :source_cell_names
	has_many :tissues, through: :source_tissue_names
	has_many :drugs, through: :source_drug_names
end
