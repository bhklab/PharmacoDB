class Tissue < ApplicationRecord
	has_many :source_tissue_names
	has_many :cell_tissues
	has_many :cells, through: :cell_tissues
	has_many :sources, through: :source_tissue_names
	has_many :experiments
	has_many :datasets, through: :experiments
	has_many :drugs, through: :experiments
	has_many :gene_drugs
end
