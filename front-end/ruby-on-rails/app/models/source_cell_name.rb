class SourceCellName < ApplicationRecord
	belongs_to :cell
	belongs_to :source
	has_many :experiments, through: :cell
end
